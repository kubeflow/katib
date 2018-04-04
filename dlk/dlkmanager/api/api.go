package api

// Learning Task Env Config
type EnvConf struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Learning Task Config
type LTConfig struct {
	PsImage     string    `json:"psImg"`     // PS image for the container in the pod(s)
	WorkerImage string    `json:"wkImg"`     // Worker image for the container in the pod(s)
	Ns          string    `json:"ns"`        // Namespace for the new pod(s)
	Scheduler   string    `json:"scheduler"` // A name of scheduler that should schedule the pod(s)
	Name        string    `json:"name"`      // A name of learning task
	NrPS        int       `json:"nrPS"`      // A number of parameter servers
	NrWorker    int       `json:"nrWorker"`  // A number of workers
	Gpu         int       `json:"gpu"`       // A number of required GPU for each task
	DryRun      bool      `json:"dry-run"`
	EntryPoint  string    `json:"entryPoint"`
	Parameters  string    `json:"parameters"`
	Timeout     int       `json:"timeout"`  // timeout of each learning task, unit is second
	Pvc         string    `json:"pvc"`      // persistent volume claim
	MountPath   string    `json:"mpath"`    // nfs mount path
	Priority    int       `json:"priority"` // learning task priority
	User        string    `json:"user"`     // user name
	Envs        []EnvConf `json:"envs"`
	PullSecret  string    `json:"pullSecret"`
}
