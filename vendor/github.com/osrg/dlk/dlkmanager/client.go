package main

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type clientConfig struct {
	Addr          string
	SchedulerName string
}

type client struct {
	apisrvClient     *kubernetes.Clientset
	unscheduledPodCh chan *api.Pod
	nodeAddCh        chan *api.Node
	nodeDeleteCh     chan *api.Node
	nodeUpdateCh     chan *api.Node
}

func checkNodeHealth(node *api.Node) bool {
	// Skip the master node
	if node.Spec.Unschedulable {
		return false
	}
	// Skip unhealty node
	for _, c := range node.Status.Conditions {
		var desireStatus string
		switch c.Type {
		case "Ready":
			desireStatus = "True"
		default:
			desireStatus = "False"
		}
		if string(c.Status) != desireStatus {
			return false
		}
	}
	for _, t := range node.Spec.Taints {
		switch t.Effect {
		case api.TaintEffectNoSchedule,
			api.TaintEffectPreferNoSchedule,
			api.TaintEffectNoExecute:
			return false
		}
	}
	return true
}
func clientNew(cfg clientConfig, podChanSize int) (*client, error) {
	//	restCfg := &restclient.Config{
	//		Host:  fmt.Sprintf("http://%s", cfg.Addr),
	//		QPS:   1000,
	//		Burst: 1000,
	//	}
	config, err := restclient.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	pch := make(chan *api.Pod, podChanSize)

	// Create informer to watch on unscheduled Pods (non-failed non-succeeded pods with an empty node binding)
	sel := "spec.nodeName==" + "" + ",status.phase!=" + string(api.PodSucceeded) + ",status.phase!=" + string(api.PodFailed)
	informer := cache.NewSharedInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.FieldSelector = sel
				return clientset.CoreV1().Pods("").List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.FieldSelector = sel
				return clientset.CoreV1().Pods("").Watch(options)
			},
		},
		&api.Pod{},
		0,
	)
	// Add event handlers for the addition, update and deletion of the pods watched by the above informer
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*api.Pod)
			// fmt.Printf("pod: %v\n", pod)
			// fmt.Printf("pod.Spec.SchedulerName: %s\n", pod.Spec.SchedulerName)
			if pod.Spec.SchedulerName != cfg.SchedulerName {
				return
			}
			pch <- pod
		},
		UpdateFunc: func(oldObj, newObj interface{}) {},
		DeleteFunc: func(obj interface{}) {},
	})
	stopCh := make(chan struct{})
	go informer.Run(stopCh)

	naddch := make(chan *api.Node, 100)
	nupdatech := make(chan *api.Node, 100)
	ndeletech := make(chan *api.Node, 100)
	// Informer for watching the addition and removal of nodes in the cluster
	lw := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "nodes", api.NamespaceAll, fields.Everything())
	_, nodeInformer := cache.NewInformer(lw, &v1.Node{}, 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				node := obj.(*api.Node)
				if !checkNodeHealth(node) {
					return
				}
				naddch <- node
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldNode := oldObj.(*api.Node)
				newNode := newObj.(*api.Node)
				if !checkNodeHealth(newNode) {
					ndeletech <- oldNode
				} else {
					nupdatech <- newNode
				}
			},
			DeleteFunc: func(obj interface{}) {
				node := obj.(*api.Node)
				ndeletech <- node
			},
		},
	)
	stopCh2 := make(chan struct{})
	go nodeInformer.Run(stopCh2)

	return &client{
		apisrvClient:     clientset,
		unscheduledPodCh: pch,
		nodeAddCh:        naddch,
		nodeDeleteCh:     ndeletech,
		nodeUpdateCh:     nupdatech,
	}, nil
}

type PodChan <-chan *api.Pod

func (c *client) GetUnscheduledPodChan() PodChan {
	return c.unscheduledPodCh
}

type NodeChan <-chan *api.Node

func (c *client) GetNodeAddChan() NodeChan {
	return c.nodeAddCh
}

func (c *client) GetNodeUpdateChan() NodeChan {
	return c.nodeUpdateCh
}

func (c *client) GetNodeDeleteChan() NodeChan {
	return c.nodeDeleteCh
}

// Write out node bindings
func (c *client) AssignBinding(ns string, bindings []*api.Binding) error {
	for _, binding := range bindings {
		err := c.apisrvClient.CoreV1().Pods(ns).Bind(binding)
		// err := c.apisrvClient.CoreV1().Pods(ns).Bind(&v1.Binding{
		// 	ObjectMeta: metav1.ObjectMeta{Namespace: ns, Name: name, UID: ob.UID},
		// 	Target: api.ObjectReference{
		// 		Kind: "Node",
		// 		Name: parseNodeID(ob.NodeID),
		// 	},
		// })
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// Returns a batch of pods or blocks until there is at least on pod creation call back
// The timeout specifies how long to wait for another pod on the pod channel before returning
// the batch of pods that need to be scheduled
func (c *client) GetPodBatch(timeout time.Duration) []*api.Pod {
	batchedPods := make([]*api.Pod, 0)

	fmt.Printf("Waiting for a pod scheduling request\n")

	// Check for first pod, block until at least 1 is available
	pod := <-c.unscheduledPodCh
	batchedPods = append(batchedPods, pod)

	// Set timer for timeout between successive pods
	timer := time.NewTimer(timeout)
	done := make(chan bool)
	go func() {
		<-timer.C
		done <- true
	}()

	fmt.Printf("Batching pod scheduling requests\n")
	numPods := 1
	//fmt.Printf("Number of pods requests: %d", numPods)
	// Poll until done from timeout
	// TODO: Put a cap on the batch size since this could go on forever
	finish := false
	for !finish {
		select {
		case pod = <-c.unscheduledPodCh:
			numPods++
			fmt.Printf("\rNumber of pods requests: %d", numPods)
			batchedPods = append(batchedPods, pod)
			// Refresh the timeout for next pod
			timer.Reset(timeout)
		case <-done:
			finish = true
			fmt.Printf("\n")
		default:
			// Do nothing and keep polling until timeout
		}
	}
	return batchedPods
}

func (c *client) GetClientset() *kubernetes.Clientset {
	return c.apisrvClient
}
