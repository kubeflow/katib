package studyjobclient

import (
	"io/ioutil"
	"strings"

	studyjobv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type StudyjobClient struct {
	config    *rest.Config
	client    *rest.RESTClient
	clientset *kubernetes.Clientset
}

func NewStudyjobClient(config *rest.Config) (*StudyjobClient, error) {
	var err error
	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	studyjobv1alpha1.AddToScheme(scheme.Scheme)
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &studyjobv1alpha1.SchemeGroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	RestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		return nil, err
	}
	return &StudyjobClient{
		config:    config,
		client:    RestClient,
		clientset: clientset,
	}, nil
}

func (s *StudyjobClient) GetStudyJobList(namespace ...string) (*studyjobv1alpha1.StudyJobList, error) {
	result := &studyjobv1alpha1.StudyJobList{}
	var ns string
	if namespace == nil {
		data, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		ns = strings.TrimSpace(string(data))
	} else {
		ns = namespace[0]
	}
	err := s.client.
		Get().
		Namespace(ns).
		Resource("studyjobs").
		Do().
		Into(result)
	return result, err
}

func (s *StudyjobClient) CreateStudyJob(studyJob *studyjobv1alpha1.StudyJob, namespace ...string) (*studyjobv1alpha1.StudyJob, error) {
	result := &studyjobv1alpha1.StudyJob{}
	var ns string
	if namespace == nil {
		data, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		ns = strings.TrimSpace(string(data))
	} else {
		ns = namespace[0]
	}
	err := s.client.
		Post().
		Namespace(ns).
		Resource("studyjobs").
		Body(studyJob).
		Do().
		Into(result)
	return result, err
}

func (s *StudyjobClient) GetWorkerTemplates() (map[string]string, error) {
	data, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	ns := strings.TrimSpace(string(data))
	cm, err := s.clientset.CoreV1().ConfigMaps(ns).Get("worker-template", metav1.GetOptions{})
	if err != nil {
		return map[string]string{}, err
	}
	return cm.Data, nil
}

func (s *StudyjobClient) UpdateWorkerTemplates(newWorkerTemplates map[string]string) error {
	data, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	ns := strings.TrimSpace(string(data))
	cm, err := s.clientset.CoreV1().ConfigMaps(ns).Get("worker-template", metav1.GetOptions{})
	if err != nil {
		return err
	}
	cm.Data = newWorkerTemplates
	_, err = s.clientset.CoreV1().ConfigMaps(ns).Update(cm)
	return err
}

func (s *StudyjobClient) GetMetricsCollectorTemplates() (map[string]string, error) {
	data, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	ns := strings.TrimSpace(string(data))
	cm, err := s.clientset.CoreV1().ConfigMaps(ns).Get("metricscollector-template", metav1.GetOptions{})
	if err != nil {
		return map[string]string{}, err
	}
	return cm.Data, nil
}

func (s *StudyjobClient) UpdateMetricsCollectorTemplates(newMCTemplates map[string]string) error {
	data, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	ns := strings.TrimSpace(string(data))
	cm, err := s.clientset.CoreV1().ConfigMaps(ns).Get("metricscollector-template", metav1.GetOptions{})
	if err != nil {
		return err
	}
	cm.Data = newMCTemplates
	_, err = s.clientset.CoreV1().ConfigMaps(ns).Update(cm)
	return err
}
