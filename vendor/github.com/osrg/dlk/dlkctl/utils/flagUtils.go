package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

//Params is struct for containing flag parameter value
type Params struct {
	Image          string
	Ns             string
	Name           string
	Scheduler      string
	NrPs           int
	NrWorker       int
	Gpu            int
	DryRun         bool
	PsRawImage     string
	WorkerRawImage string
	BaseImage      string
	EntryPoint     string
	Parameters     string
	Timeout        int
	Pvc            string
	MountPath      string
	Priority       int
	SinceTime      string
	GpuImg         bool
}

const (
	PrioMIN = 0   // lowest priority
	PrioMAX = 100 // highest priority
)

//CheckFlags checks flags value vailidity and return these value
func CheckFlags(cmd *cobra.Command) (p Params, er error) {
	var rtn = Params{}
	var err error
	e := false

	//image parameter
	if cmd.Flags().Lookup("image") != nil {
		rtn.Image, err = cmd.Flags().GetString("image")
		if err != nil {
			return Params{}, err
		}
	}

	//namespace parameter
	if cmd.Flags().Lookup("ns") != nil {
		rtn.Ns, err = cmd.Flags().GetString("ns")
		if err != nil {
			return Params{}, err
		}
	}

	//learning task name parameter
	if cmd.Flags().Lookup("name") != nil {
		rtn.Name, err = cmd.Flags().GetString("name")
		if err != nil {
			return Params{}, err
		}
	}

	//scheduler parameter
	if cmd.Flags().Lookup("scheduler") != nil {
		rtn.Scheduler, err = cmd.Flags().GetString("scheduler")
		if err != nil {
			return Params{}, err
		}
	}

	//number of parameter server
	if cmd.Flags().Lookup("nrPs") != nil {
		rtn.NrPs, err = cmd.Flags().GetInt("nrPs")
		if err != nil {
			return Params{}, err
		} else if rtn.NrPs < 0 {
			fmt.Println("flag: --nrPs must be equal to or greater than 0")
			e = true
		}
	}

	//number of worker parameter
	if cmd.Flags().Lookup("nrWorker") != nil {
		rtn.NrWorker, err = cmd.Flags().GetInt("nrWorker")
		if err != nil {
			return Params{}, err
		} else if rtn.NrWorker < 1 {
			fmt.Println("flag: --nrWorker must be greater than 0")
			e = true
		}
	}

	// number of worker must be 1
	// in case number of parameter server is 0
	if rtn.NrPs == 0 && rtn.NrWorker > 1 {
		fmt.Println("flag: --nrWorker must be 1 in case --nrPs is 0")
		e = true
	}

	//Gpu parameter
	if cmd.Flags().Lookup("gpu") != nil {
		rtn.Gpu, err = cmd.Flags().GetInt("gpu")
		if err != nil {
			return Params{}, err
		} else if rtn.Gpu < 0 {
			fmt.Println("flag: --gpu must be equal to or greater than 0")
			e = true
		}
	}

	// timeout parameter
	if cmd.Flags().Lookup("timeout") != nil {
		rtn.Timeout, err = cmd.Flags().GetInt("timeout")
		if err != nil {
			return Params{}, err
		}
	}

	// persistent volume claim parameter
	if cmd.Flags().Lookup("pvc") != nil {
		rtn.Pvc, err = cmd.Flags().GetString("pvc")
		if err != nil {
			return Params{}, err
		}
	}

	// nfs mount path parameter
	if cmd.Flags().Lookup("mpath") != nil {
		rtn.MountPath, err = cmd.Flags().GetString("mpath")
		if err != nil {
			return Params{}, err
		}
		// Check whether the path is absolute
		if rtn.MountPath != "" && !filepath.IsAbs(rtn.MountPath) {
			fmt.Println("flag: --mpath specified path is not absolute")
			e = true
		}
	}

	// Without persistent volume claim, nfs mount path is specified only
	if rtn.Pvc == "" && rtn.MountPath != "" {
		fmt.Println("flag: --mpath must be specified along with persistent volume claim")
		e = true
	}

	// priority parameter
	if cmd.Flags().Lookup("priority") != nil {
		rtn.Priority, err = cmd.Flags().GetInt("priority")
		if err != nil {
			return Params{}, err
		} else if rtn.Priority < PrioMIN || rtn.Priority > PrioMAX {
			fmt.Println("flag: --priority must be from 0(lowest) to 100(highest)")
			e = true
		}
	}

	//ps raw image parameter
	if cmd.Flags().Lookup("psRawImage") != nil {
		rtn.PsRawImage, err = cmd.Flags().GetString("psRawImage")
		if err != nil {
			return Params{}, err
		}
	}

	//worker raw image parameter
	if cmd.Flags().Lookup("workerRawImage") != nil {
		rtn.WorkerRawImage, err = cmd.Flags().GetString("workerRawImage")
		if err != nil {
			return Params{}, err
		}
	}

	//both image and raw images can't be specified at the same time
	if rtn.Image != "" && (rtn.PsRawImage != "" || rtn.WorkerRawImage != "") {
		fmt.Println("don't specify both --image and --ps/workerRawImage")
		e = true
	}

	//docker base image parameter
	if cmd.Flags().Lookup("baseImage") != nil {
		rtn.BaseImage, err = cmd.Flags().GetString("baseImage")
		if err != nil {
			return Params{}, err
		}
	}

	//docker container entry point
	if cmd.Flags().Lookup("entryPoint") != nil {
		rtn.EntryPoint, err = cmd.Flags().GetString("entryPoint")
		if err != nil {
			return Params{}, err
		}
	}

	//docker container exec parameters
	if cmd.Flags().Lookup("parameters") != nil {
		rtn.Parameters, err = cmd.Flags().GetString("parameters")
		if err != nil {
			return Params{}, err
		}
	}

	//gpu image parameter
	if cmd.Flags().Lookup("gpu-image") != nil {
		rtn.GpuImg, err = cmd.Flags().GetBool("gpu-image")
		if err != nil {
			return Params{}, err
		}
	}

	//DryRun parameter
	if cmd.Flags().Lookup("dry-run") != nil {
		rtn.DryRun, err = cmd.Flags().GetBool("dry-run")
		if err != nil {
			return Params{}, err
		} else if e {
			cmd.Help()
			os.Exit(1)
		}
	}

	//start time of duration specification
	if cmd.Flags().Lookup("sinceTime") != nil {
		rtn.SinceTime, err = cmd.Flags().GetString("sinceTime")
		//Require RFC3339 formated string
		if rtn.SinceTime != "" {
			_, err := time.Parse(time.RFC3339, rtn.SinceTime)
			if err != nil {
				return Params{}, err
			}
		}
	}

	return rtn, err
}

// AddDryRunFlag set dry run flag to passed command
func AddDryRunFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("dry-run", false, "only print the object that would be sent,without sending it")
}

// AddImageFlag set docker image name for learning task
func AddImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("image", "", "set docker image name")
}

// AddNameSpaceFlag identify which namespace in k8s request will be deployed
func AddNameSpaceFlag(cmd *cobra.Command) {
	cmd.Flags().String("ns", "default", "set namespace in which request deploy")
}

// AddNameFlag set learning task name
func AddNameFlag(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "set learning task name,if not present,dlkctl generate it automatically")
}

// AddSchedulerFlag set scheduler name which pick up and assign the dlk request to node
func AddSchedulerFlag(cmd *cobra.Command) {
	cmd.Flags().String("scheduler", "dlk", "set scheduler name which pick up and assign the dlk request to node")
}

// AddNrPsFlag set number of parameter server
func AddNrPsFlag(cmd *cobra.Command) {
	cmd.Flags().Int("nrPs", 1, "set number of PS")
}

// AddNrWorkerFlag set number of worker server
func AddNrWorkerFlag(cmd *cobra.Command) {
	cmd.Flags().Int("nrWorker", 1, "set number of Worker")
}

//AddGpuFlag enable use of GPU tensorflow Image and set limit number of gpu
func AddGpuFlag(cmd *cobra.Command) {
	cmd.Flags().Int("gpu", 0, "use gpu tensorflow image and set limit number of gpu ")
}

// AddPsRawImageFlag set ps docker image name (directly passed to job objects)
func AddPsRawImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("psRawImage", "", "set ps docker image name (directly passed to job objects)")
}

// AddWorkerRawImageFlag set worker docker image name (directly passed to job objects)
func AddWorkerRawImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("workerRawImage", "", "set worker docker image name (directly passed to job objects)")
}

//AddTypeFlag add a flag which filter output by their role
func AddTypeFlag(cmd *cobra.Command) {
	cmd.Flags().String("type", "", "expected value: [worker,ps]. display only worker or ps related object")
}

//AddBaseImageFlag add a flag which specify can docker base image
func AddBaseImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("baseImage", "", "docker base image")
}

//AddEntryPointFlag add a flag which specify docker base image
func AddEntryPointFlag(cmd *cobra.Command) {
	cmd.Flags().String("entryPoint", "", "docker container entry point")
}

//AddParametersFlag add a flag which specify docker container cmd parameter
func AddParametersFlag(cmd *cobra.Command) {
	cmd.Flags().String("parameters", "", "docker container exec parameter")
}

//AddTimeoutFlag add a flag which specify timeout of each learning task
func AddTimeoutFlag(cmd *cobra.Command) {
	cmd.Flags().Int("timeout", 0, "timeout of each learning task")
}

//AddPvcFlag add a flag which specify persistent volume claim
func AddPvcFlag(cmd *cobra.Command) {
	cmd.Flags().String("pvc", "", "persistent volume claim")
}

//AddMountPathFlag add a flag which specify nfs mount path
func AddMountPathFlag(cmd *cobra.Command) {
	cmd.Flags().String("mpath", "", "nfs mount path (dafault \"/default-path\")")
}

//AddPriorityFlag add a flag which specify learning task priority
func AddPriorityFlag(cmd *cobra.Command) {
	cmd.Flags().Int("priority", 0, "learning task priority (default 0 | 0:lowest - 100:highest)")
}

//AddSinceTimeFlag add a flag which specify start time of duration specification
func AddSinceTimeFlag(cmd *cobra.Command) {
	cmd.Flags().String("sinceTime", "", "Only return logs after a specific time (RFC3339).defaults to all logs")
}

//AddGpuImageFlag set a flag which specify whether gpu image is created or not
func AddGpuImageFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("gpu-image", false, "create docker image for gpu as well (dafault false)")
}
