package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/osrg/dlk/dlkmanager/configs"

	"github.com/osrg/dlk/dlkmanager/api"
	"github.com/osrg/dlk/dlkmanager/datastore"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"

	"github.com/labstack/echo" // Web server framework for REST API
)

const (
	GPUMIN  = 0   // minimum GPU number
	GPUMAX  = 3   // maximum GPU number
	PrioMIN = 0   // lowest learning task priority
	PrioMAX = 100 // highest learning task priority
)

// POST's JSON
type RunParam struct {
	LtConf *api.LTConfig
}

func (r *RunParam) runParamChk() (code int) {

	// check whether both ps and worker images are specified
	if r.LtConf.PsImage == "" || r.LtConf.WorkerImage == "" {
		log.Error("ps image and/or worker image are not specified")
		return http.StatusBadRequest
	}

	if r.LtConf.Ns == "" {
		r.LtConf.Ns = configs.Pflg.Ns
	}

	if r.LtConf.Scheduler == "" {
		r.LtConf.Scheduler = configs.Pflg.Scheduler
	}

	if r.LtConf.Name == "" {
		r.LtConf.Name = ""
	}

	if r.LtConf.NrPS < 0 {
		log.Error("NrPS is less than 0")
		return http.StatusBadRequest
	}

	if r.LtConf.NrWorker < 1 {
		log.Error("NrWorker is less than 1")
		return http.StatusBadRequest
	}

	// number of worker must be 1
	// in case number of parameter server is 0
	if r.LtConf.NrPS == 0 && r.LtConf.NrWorker > 1 {
		log.Error("NrWorker must be 1 in case NrPs is 0")
		return http.StatusBadRequest
	}

	if r.LtConf.Gpu < GPUMIN || r.LtConf.Gpu > GPUMAX {
		log.Error("Gpu is out of range.")
		return http.StatusBadRequest
	}

	if r.LtConf.Priority < PrioMIN || r.LtConf.Priority > PrioMAX {
		log.Error("Priority is out of range.")
		return http.StatusBadRequest
	}

	if r.LtConf.EntryPoint == "" {
		r.LtConf.EntryPoint = "python /distributed-tensorflow-example.py"
	}

	if r.LtConf.MountPath == "" {
		r.LtConf.MountPath = "/default-path"
	}

	// check mount path is absolute
	if !filepath.IsAbs(r.LtConf.MountPath) {
		log.Error("mount path is not absolute")
		return http.StatusBadRequest
	}

	log.Info(fmt.Sprintf("After RunParam is checked, RunParam = %#v", r))

	return http.StatusOK
}

func (r *RunParam) learningTaskInit() (lt *learningTask, code int) {
	// Initialize the kubernetes client
	//	restCfg := &restclient.Config{
	//		Host:  fmt.Sprintf("http://%s", configs.Pflg.Addr),
	//		QPS:   1000,
	//		Burst: 1000,
	//	}
	// Create the kubernetes client

	config, err := restclient.InClusterConfig()
	if err != nil {
		return nil, http.StatusBadRequest
	}
	k8scli, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(fmt.Sprintf("kubernetes client not created. error = %d", err))
		return nil, http.StatusInternalServerError
	}

	// Check PersistentVolumeClaim
	var pvc *corev1.PersistentVolumeClaim
	var e bool
	if r.LtConf.Pvc != "" {
		pvc, e = r.checkPersistentVolumeClaim(k8scli)
		if e {
			return nil, http.StatusBadRequest
		}
	}

	if r.LtConf.Name == "" {
		now := time.Now()
		r.LtConf.Name = fmt.Sprintf("dlk-%d-%d-%d-%d-%d-%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		log.Info(fmt.Printf("generated learning task name: %s", r.LtConf.Name))
	}

	lt = newLearningTask(r.LtConf, k8scli, pvc)
	return lt, http.StatusOK
}

// Check PersistentVolumeClaim
func (r *RunParam) checkPersistentVolumeClaim(k8scli *kubernetes.Clientset) (pvc *corev1.PersistentVolumeClaim, e bool) {

	// Get PersistentVolumeClaim
	getopt := metav1.GetOptions{}
	getpvc, err := k8scli.CoreV1().PersistentVolumeClaims(r.LtConf.Ns).Get(r.LtConf.Pvc, getopt)
	if err != nil {
		log.Error(fmt.Sprintf("Specified pvc does not exist in combination with the specified namespace. error = %d", err))
		return nil, true
	}

	ret := false
	// Check Bind
	if getpvc.Status.Phase != corev1.ClaimBound {
		log.Error("pvc is not bound")
		ret = true
	}

	// Check AccessModes
	incl := false
	for _, mode := range getpvc.Status.AccessModes {
		if mode == corev1.ReadWriteMany {
			incl = true
			break
		}
	}
	if !incl {
		log.Error("pvc access modes don't include ReadWriteMany")
		ret = true
	}

	return getpvc, ret
}

// POST handling
func runLearningTask(c echo.Context) error {

	var code int

	// Extracts JSON
	r := RunParam{LtConf: &api.LTConfig{}}
	err := c.Bind(r.LtConf)
	if err != nil {
		log.Error(fmt.Sprintf("POST: JSON is not extracted. error = %d", err))
		return c.String(http.StatusInternalServerError, "Internal error.")
	}

	log.Info(fmt.Sprintf("POST: JSON is extracted. RunParam = %#v", r))

	code = r.runParamChk()
	if code != http.StatusOK {
		return c.String(code, "Input parameter value(s) are not correct.")
	}

	lt, code := r.learningTaskInit()
	if code != http.StatusOK {
		return c.String(code, "Internal error.")
	}

	if r.LtConf.DryRun {
		return c.String(http.StatusOK, "Dry run is performed.")
	}

	// Respond to CLI immediately, before learning task is finished.
	c.String(http.StatusOK, "Command is performed.")

	//store learning task information into datastore for check learning tasks from get method
	info := datastore.LearningTaskInfo{
		Name:        lt.ltc.Name,
		PsImage:     lt.ltc.PsImage,
		WorkerImage: lt.ltc.WorkerImage,
		Ns:          lt.ltc.Ns,
		Scheduler:   lt.ltc.Scheduler,
		NrPS:        lt.ltc.NrPS,
		NrWorker:    lt.ltc.NrWorker,
		Gpu:         lt.ltc.Gpu,
		Created:     time.Now(),
		Timeout:     lt.ltc.Timeout,
		Pvc:         lt.ltc.Pvc,
		MountPath:   lt.ltc.MountPath,
		Priority:    lt.ltc.Priority,
		State:       ltStateNotCompleted,
		User:        lt.ltc.User,
		PodState:    make(map[string]string),
	}
	datastore.Accesor.Put(info)

	go lt.run()

	return err
}
