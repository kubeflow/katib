package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kubeflow/katib/dlk/dlkmanager/api"
	"github.com/kubeflow/katib/dlk/dlkmanager/datastore"

	lgr "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ltStateNotCompleted = "not completed"
	ltStateCompleted    = "completed"
	ltStateStopped      = "stopped"
	ltStateDeleted      = "deleted"
	ltStateTimeout      = "timeout"
	ltStateRunning      = "running"
)

// Learning Task
type learningTask struct {
	c    *kubernetes.Clientset
	ltc  *api.LTConfig
	pvc  *apiv1.PersistentVolumeClaim
	name string

	workers            map[string]struct{}
	nrCompletedWorkers int
	workerCmpCh        chan string

	psSvcs     []*psService
	workerSvcs []*workerService

	psJobs     []*psJob
	workerJobs []*workerJob

	deleteCh chan bool
	stopCh   chan (chan struct{})

	psLogs     map[string][]logObj
	workerLogs map[string][]logObj

	nrWorkers, nrPSes           int
	nrReadyWorkers, nrReadyPSes int
	pods                        []*apiv1.Pod

	running bool

	usingGPUsPerPod map[*apiv1.Pod]int
	podToNode       map[*apiv1.Pod]*apiv1.Node
}

// Log Object
type logObj struct {
	time  metav1.Time
	value string
}

// Create Services
func (lt *learningTask) createServices() {
	lt.psSvcs = newPSServices(lt.name, lt.ltc.NrPS)
	lt.workerSvcs = newWorkerServices(lt.name, lt.ltc.NrWorker)

	for _, svc := range lt.psSvcs {
		_, err := lt.c.CoreV1().Services(lt.ltc.Ns).Create(svc.svc)
		if err != nil {
			log.Error(fmt.Sprintf("failed to create a PS service %s: %s", svc.name, err))
		}
	}

	for _, svc := range lt.workerSvcs {
		_, err := lt.c.CoreV1().Services(lt.ltc.Ns).Create(svc.svc)
		if err != nil {
			log.Error(fmt.Sprintf("failed to create a worker service %s: %s", svc.name, err))
		}
	}
}

// Delete Services
func (lt *learningTask) deleteServices() {
	for _, svc := range lt.psSvcs {
		err := lt.c.CoreV1().Services(lt.ltc.Ns).Delete(svc.name, &metav1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("failed to delete a PS service %s: %s", svc.name, err))
		}
	}

	for _, svc := range lt.workerSvcs {
		err := lt.c.CoreV1().Services(lt.ltc.Ns).Delete(svc.name, &metav1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("failed to delete a worker service %s: %s", svc.name, err))
		}
	}
}

// Generate Member Lists
func (lt *learningTask) genMemberLists() []string {
	workers := "--workers="
	for i, workerSvc := range lt.workerSvcs {
		if i == 0 {
			workers += fmt.Sprintf("%s:2222", workerSvc.name)
		} else {
			workers += fmt.Sprintf(",%s:2222", workerSvc.name)
		}
	}

	pses := "--parameter_servers="
	for i, psSvc := range lt.psSvcs {
		if i == 0 {
			pses += fmt.Sprintf("%s:2222", psSvc.name)
		} else {
			pses += fmt.Sprintf(",%s:2222", psSvc.name)
		}
	}

	return []string{
		workers,
		pses,
	}
}

// Create Jobs
func (lt *learningTask) createJobs() {
	memberLists := lt.genMemberLists()

	lt.psJobs = newPSJobs(lt, memberLists)
	lt.workerJobs = newWorkerJobs(lt, memberLists)

	// Dry Run
	if lt.ltc.DryRun {
		log.Info(fmt.Sprintf("Dry Run: PS jobs = %#v", lt.psJobs))
		log.Info(fmt.Sprintf("Dry Run: Worker jobs = %#v", lt.workerJobs))
		return
	}

	for _, job := range lt.psJobs {
		_, err := lt.c.BatchV1().Jobs(lt.ltc.Ns).Create(job.job)
		if err != nil {
			log.Error(fmt.Sprintf("failed to create a PS job %s: %s", job.name, err))
		}
		lt.psLogs[job.name] = []logObj{}
	}

	for _, job := range lt.workerJobs {
		_, err := lt.c.BatchV1().Jobs(lt.ltc.Ns).Create(job.job)
		if err != nil {
			log.Error(fmt.Sprintf("failed to create a worker job %s: %s", job.name, err))
		}

		lt.workers[job.name] = struct{}{}
		lt.workerLogs[job.name] = []logObj{}
	}
}

// Cleanup Jobs
func (lt *learningTask) cleanupJobs() {
	for _, job := range lt.psJobs {
		err := lt.c.BatchV1().Jobs(lt.ltc.Ns).Delete(job.name, &metav1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("failed to delete a PS job %s: %s", job.name, err))
		}
	}

	for _, job := range lt.workerJobs {
		err := lt.c.BatchV1().Jobs(lt.ltc.Ns).Delete(job.name, &metav1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("failed to delete a worker job %s: %s", job.name, err))
		}
	}

	pods, err := lt.c.CoreV1().Pods(lt.ltc.Ns).List(metav1.ListOptions{LabelSelector: fmt.Sprintf("learning-task=%s", lt.name)})
	if err != nil {
		log.Error(fmt.Sprintf("failed to delete pods: %s", err))
	}

	for _, pod := range pods.Items {
		log.Info(fmt.Sprintf("deleting pod %s", pod.ObjectMeta.Name))
		err = lt.c.CoreV1().Pods(lt.ltc.Ns).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("failed to delete pod %s: %s", pod.ObjectMeta.Name, err))
		}
	}
}

// Get Logs
func (lt *learningTask) getLogs(jobname string, namespace string, Logs *map[string][]logObj) {
	pl, err := lt.c.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: "job-name=" + jobname})
	if err != nil {
		log.Error(fmt.Sprintf("failed to obtain pod list: %s", err))
		os.Exit(1)
	}
	if len(pl.Items) != 0 {
		var logs []byte
		if len((*Logs)[jobname]) == 0 {
			logs, _ = lt.c.CoreV1().Pods(namespace).GetLogs(pl.Items[0].ObjectMeta.Name, &apiv1.PodLogOptions{Timestamps: true}).Do().Raw()
		} else {
			logs, _ = lt.c.CoreV1().Pods(namespace).GetLogs(pl.Items[0].ObjectMeta.Name, &apiv1.PodLogOptions{SinceTime: &(*Logs)[jobname][len((*Logs)[jobname])-1].time, Timestamps: true}).Do().Raw()
		}
		if len(logs) > 1 && pl.Items[0].Status.Phase != apiv1.PodPending && pl.Items[0].Status.Phase != apiv1.PodUnknown {
			logf := strings.Split(string(logs), "\n")
			for _, l := range logf {
				ll := strings.SplitN(l, " ", 2)
				if len(ll) == 2 {
					t, err := time.Parse(time.RFC3339, ll[0])
					if err == nil && (len((*Logs)[jobname]) == 0 || (*Logs)[jobname][len((*Logs)[jobname])-1].time.Time.Before(t)) {
						(*Logs)[jobname] = append((*Logs)[jobname], logObj{time: metav1.Time{Time: t}, value: ll[1]})
					}
				}
			}
		}
	}
}

// Polling Jobs
func (lt *learningTask) pollJobs() {
	for _, job := range lt.psJobs {
		lt.getLogs(job.name, lt.ltc.Ns, &lt.psLogs)
	}
	for _, job := range lt.workerJobs {
		lt.getLogs(job.name, lt.ltc.Ns, &lt.workerLogs)
	}

	if lt.nrCompletedWorkers == lt.ltc.NrWorker {
		return
	}

	jobList, err := lt.c.BatchV1().Jobs(lt.ltc.Ns).List(metav1.ListOptions{}) // TODO: label selector
	if err != nil {
		log.Error(fmt.Sprintf("failed to obtain job list: %s", err))
		os.Exit(1)
	}

	for _, job := range jobList.Items {
		if len(job.Status.Conditions) == 0 {
			continue
		}

		cond := job.Status.Conditions[0]
		if _, ok := lt.workers[job.ObjectMeta.Name]; ok {
			if cond.Type == batchv1.JobComplete || cond.Type == batchv1.JobFailed {
				go func(lt *learningTask, name string) {
					lt.workerCmpCh <- name
				}(lt, job.ObjectMeta.Name)
			}
		}
		delete(lt.workers, job.ObjectMeta.Name)
	}

}

// Run Learning Task
func (lt *learningTask) run() {

	lt.createServices()
	lt.createJobs()

	running := true

	timeoutCh := make(<-chan time.Time)
	if lt.ltc.Timeout != 0 {
		timeoutCh = time.After(time.Duration(lt.ltc.Timeout) * time.Second)
	}

	// local learningtask status
	state := ltStateNotCompleted

	// local pods status
	podState := make(map[string]apiv1.PodPhase)
	for _, worker := range lt.workerJobs {
		podState[worker.name] = apiv1.PodPending
	}
	for _, ps := range lt.psJobs {
		podState[ps.name] = apiv1.PodPending
	}
	// set init pod status value to datastore
	for k := range podState {
		datastore.Accesor.UpdatePodState(lt.name, k, ltStateNotCompleted)
	}

	var stopCompleteNotifyCh chan struct{}

	for running {
		select {
		case cmpedWorker := <-lt.workerCmpCh:
			log.Info(fmt.Sprintf("worker %s completed", cmpedWorker))
			lt.nrCompletedWorkers++

			// if it is first complete worker,then set status to not complete
			if lt.nrCompletedWorkers == 1 {
				state = ltStateNotCompleted
				datastore.Accesor.UpdateState(lt.name, state, "")
			}

			if lt.nrCompletedWorkers == lt.ltc.NrWorker {
				state = ltStateCompleted
				running = false
				break
			}

		case <-time.After(1 * time.Second):
			err := lt.checkPodStatus(podState)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			lt.pollJobs()

			if lt.nrCompletedWorkers == lt.ltc.NrWorker {
				state = ltStateCompleted
				running = false
				break
			}

			// if all worker is running,set state to running
			if lt.nrCompletedWorkers == 0 && state != ltStateRunning {
				i := 0
				for _, pState := range podState {
					if pState != apiv1.PodRunning {
						break
					}
					i++
					if i == len(podState) {
						state = ltStateRunning
						datastore.Accesor.UpdateState(lt.name, state, "")
					}
				}
			}

		case stopCompleteNotifyCh = <-lt.stopCh:
			log.Info(fmt.Sprintf("learning task %s stopped", lt.name))
			state = ltStateStopped
			running = false

		case <-lt.deleteCh:
			log.Info(fmt.Sprintf("learning task %s deleted", lt.name))
			state = ltStateDeleted
			running = false

		case <-timeoutCh:
			log.Infof("learning task %s stopped because of timeout (%d sec)", lt.name, lt.ltc.Timeout)
			state = ltStateTimeout
			running = false
		}
	}

	lt.cleanupJobs()
	lt.deleteServices()

	// update PS pod state since it already deleted
	for _, ps := range lt.psJobs {
		datastore.Accesor.UpdatePodState(lt.name, ps.name, ltStateCompleted)
	}

	// learning task exec time is calculated after task completion
	et := ""
	if state == ltStateCompleted {
		// get learning task created time
		ct, _ := datastore.Accesor.Get(lt.name)
		dura := time.Since(ct.Created)
		et = dura.Truncate(time.Millisecond).String()
	}

	datastore.Accesor.UpdateState(lt.name, state, et)

	if state == ltStateStopped {
		stopCompleteNotifyCh <- struct{}{}
		return
	}

	completedLTMu.Lock()
	completedLearningTasks[lt.name] = runningLearningTasks[lt.name]
	completedLTMu.Unlock()

	runningLTMu.Lock()
	delete(runningLearningTasks, lt.name)
	runningLTMu.Unlock()

	log.WithFields(
		lgr.Fields{
			"learningTask": lt.name,
			"state":        state,
		}).Info("learning task is completed")
}

var (
	runningLearningTasks   map[string]*learningTask
	completedLearningTasks map[string]*learningTask
	runningLTMu            sync.Mutex
	completedLTMu          sync.Mutex
)

func init() {
	runningLearningTasks = make(map[string]*learningTask)
	completedLearningTasks = make(map[string]*learningTask)
}

// Creates Learning Task
func newLearningTask(lc *api.LTConfig, c *kubernetes.Clientset, pvc *apiv1.PersistentVolumeClaim) *learningTask {
	ret := &learningTask{
		c:                  c,
		ltc:                lc,
		pvc:                pvc,
		name:               lc.Name,
		nrCompletedWorkers: 0,
		workerCmpCh:        make(chan string),
		workers:            make(map[string]struct{}),

		deleteCh:   make(chan bool),
		stopCh:     make(chan (chan struct{})),
		psLogs:     make(map[string][]logObj),
		workerLogs: make(map[string][]logObj),

		nrPSes:    lc.NrPS,
		nrWorkers: lc.NrWorker,

		pods: make([]*apiv1.Pod, 0),

		usingGPUsPerPod: make(map[*apiv1.Pod]int),
		podToNode:       make(map[*apiv1.Pod]*apiv1.Node),
	}

	runningLTMu.Lock()
	runningLearningTasks[lc.Name] = ret
	runningLTMu.Unlock()

	return ret
}

func (lt *learningTask) checkPodStatus(podState map[string]apiv1.PodPhase) error {
	// get pod infomation using learning-task label
	label := fmt.Sprintf("learning-task=%s", lt.name)
	pods, err := lt.c.Core().Pods("").List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return err
	} else if pods.Size() == 0 {
		return fmt.Errorf("no  pod related to  \"%s\" is found", lt.name)
	}

	for _, pod := range pods.Items {
		job := pod.Labels["job-name"]
		current, ok := podState[job]
		if !ok {
			// must be dead code
			log.Warnf("unknown job is detected. job name: %s", job)
		}

		// if pod status change,then update local status
		if pod.Status.Phase == current || pod.Status.Phase == apiv1.PodPending {
			continue
		}

		podState[job] = pod.Status.Phase

		// update pod status
		if pod.Status.Phase == apiv1.PodRunning {
			// PodRunning == "running"
			datastore.Accesor.UpdatePodState(lt.name, job, ltStateRunning)
		} else if pod.Status.Phase == apiv1.PodSucceeded {
			// PodSuccessed == "complete"
			datastore.Accesor.UpdatePodState(lt.name, job, ltStateCompleted)
		} else { // Init,pending,containerCreating,Error = "not complete"
			datastore.Accesor.UpdatePodState(lt.name, job, ltStateNotCompleted)
		}
	}

	return nil
}
