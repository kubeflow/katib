package datastore

import (
	"time"
)

// Log Object
type LogObj struct {
	Time  string `json"time"`
	Value string `json"value"`
}

// Pod Log Info
type PodLogInfo struct {
	PodName string   `json"podname"`
	Logs    []LogObj `json:"logs"`
}

// LearningTask Log Info
type LtLogInfo struct {
	LtName  string       `json:"ltname"`
	PodLogs []PodLogInfo `json:"podlogs"`
}

// LearningTaskInfo is structure for experiment and related job,svc list
type LearningTaskInfo struct {
	PsImage     string   `json:"psImg"`
	WorkerImage string   `json:"wkImg"`
	Ns          string   `json:"ns"`
	Scheduler   string   `json:"scheduler"`
	Name        string   `json:"name"`
	NrPS        int      `json:"nrPS"`
	NrWorker    int      `json:"nrWorker"`
	Gpu         int      `json:"gpu"`
	Jobs        []string `json:"jobs"`
	Services    []string `json:"services"`
	Created     time.Time
	ExecTime    string
	User        string            `json:"user"`
	Timeout     int               `json:"timeout"`
	Pvc         string            `json:"pvc"`
	MountPath   string            `json:"mpath"`
	Priority    int               `json:"priority"`
	State       string            `json:"state"`
	PodState    map[string]string `json:"podState"`
}

// LearningTaskData is interface for manage requested experiments
type LearningTaskData interface {
	Get(string) (LearningTaskInfo, error)
	GetAll() ([]LearningTaskInfo, error)
	Put(LearningTaskInfo) error
	Remove(string) error
	UpdateState(lt string, state string, time string) error
	UpdatePodState(lt string, pod string, state string) error
}

//Accesor is access interface for experiment datastore
var Accesor LearningTaskData

func init() {
	Accesor = GetLearningTaskMap()
}
