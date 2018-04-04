package nvdocker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dclient "github.com/docker/docker/client"
	"github.com/mlkube/katib/api"
	"github.com/mlkube/katib/db"
	"github.com/mlkube/katib/manager/modeldb"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	nvCtlDevice      string = "/dev/nvidiactl"
	nvUVMDevice      string = "/dev/nvidia-uvm"
	nvUVMToolsDevice string = "/dev/nvidia-uvm-tools"
	devDirectory            = "/dev"
)

var (
	nvidiaDeviceRE   = regexp.MustCompile(`^nvidia[0-9]*$`)
	nvidiaFullpathRE = regexp.MustCompile(`^/dev/nvidia[0-9]*$`)
)

type nvGPUManager struct {
	mux sync.Mutex
	// All gpus available on the Node
	allGPUs             []string
	gPUAllocedContainer map[string]string
}

func NewNvGPUManager() (*nvGPUManager, error) {
	n := &nvGPUManager{
		allGPUs:             []string{},
		gPUAllocedContainer: map[string]string{},
	}
	err := n.discoverGPUs()
	return n, err
}

func (ngm *nvGPUManager) discoverGPUs() error {
	files, err := ioutil.ReadDir(devDirectory)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if nvidiaDeviceRE.MatchString(f.Name()) {
			fmt.Printf("Found Nvidia GPU %v\n", f.Name())
			gn := path.Join(devDirectory, f.Name())
			ngm.allGPUs = append(ngm.allGPUs, gn)
			ngm.gPUAllocedContainer[gn] = ""
		}
	}
	return nil
}

func (ngm *nvGPUManager) AllocGPU(num int, cid string) (bool, []string, error) {
	ngm.mux.Lock()
	defer ngm.mux.Unlock()
	fGpusId := []int{}
	for i, g := range ngm.allGPUs {
		if ngm.gPUAllocedContainer[g] == "" {
			fGpusId = append(fGpusId, i)
		}
	}
	if len(fGpusId) < num {
		return false, nil, nil
	}
	fGpusId = fGpusId[:num]
	retid := make([]string, num)
	for j, i := range fGpusId {
		ngm.gPUAllocedContainer[ngm.allGPUs[i]] = cid
		retid[j] = strconv.Itoa(i)
	}
	log.Printf("%v GPU Alloced. ID %v", num, retid)
	return true, retid, nil
}

func (ngm *nvGPUManager) ReleaseGPU(cid string) error {
	ngm.mux.Lock()
	defer ngm.mux.Unlock()
	for _, g := range ngm.allGPUs {
		if ngm.gPUAllocedContainer[g] == cid {
			log.Printf("GPU ID %v Releaseced.", g)
			ngm.gPUAllocedContainer[g] = ""
		}
	}
	return nil
}

func (ngm *nvGPUManager) GetAllGPU() []string {
	return ngm.allGPUs
}

type spawnReq struct {
	StudyID string
	Trials  []*api.Trial
	CConf   []*container.Config
}

type NvDockerWorkerInterface struct {
	PendingTrialList   map[string][]*api.Trial
	RunningTrialList   map[string][]*api.Trial
	CompletedTrialList map[string][]*api.Trial
	dcli               *dclient.Client
	mux                *sync.Mutex
	dbIf               db.VizierDBInterface
	tidToCid           map[string]string
	ngm                *nvGPUManager
	addTrialQueCh      chan spawnReq
	deleteTrialQueCh   chan string
	stopSchedule       chan bool
}

func NewNvDockerWorkerInterface() *NvDockerWorkerInterface {
	dc, err := dclient.NewEnvClient()
	if err != nil {
		log.Printf("docker client err %v", err)
		return nil
	}
	ngm, err := NewNvGPUManager()
	if err != nil {
		log.Printf("GPU device detect err %v", err)
		return nil
	}
	n := &NvDockerWorkerInterface{
		PendingTrialList:   make(map[string][]*api.Trial),
		RunningTrialList:   make(map[string][]*api.Trial),
		CompletedTrialList: make(map[string][]*api.Trial),
		dcli:               dc,
		mux:                new(sync.Mutex),
		dbIf:               db.New(),
		ngm:                ngm,
		tidToCid:           make(map[string]string),
		addTrialQueCh:      make(chan spawnReq),
		deleteTrialQueCh:   make(chan string),
		stopSchedule:       make(chan bool),
	}
	go n.schedulingLoop()
	return n
}

func (n *NvDockerWorkerInterface) getcid(tID string) (string, error) {
	cid, ok := n.tidToCid[tID]
	if !ok {
		return "", errors.New(fmt.Sprintf("No container TID %v", tID))
	}
	return cid, nil
}

func (n *NvDockerWorkerInterface) getCon(tID string) (types.ContainerJSON, error) {
	cid, err := n.getcid(tID)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return n.dcli.ContainerInspect(context.Background(), cid)
}

func (n *NvDockerWorkerInterface) getConLog(tID string, since string) ([]string, error) {
	cid, err := n.getcid(tID)
	if err != nil {
		log.Printf("getcid err %v", err)
		return nil, err
	}
	out, err := n.dcli.ContainerLogs(context.Background(), cid, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Since: since, Timestamps: true})
	if err != nil {
		log.Printf("get log err %v", err)
		return nil, err
	}
	defer out.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	return strings.Split(buf.String(), "\n"), nil
}

func (n *NvDockerWorkerInterface) IsTrialComplete(studyId string, tID string) (bool, error) {
	c, err := n.getCon(tID)
	if err != nil {
		return false, err
	}
	if !c.State.Running && !c.State.Dead && !c.State.OOMKilled {
		return true, nil
	}
	return false, nil
}

func (n *NvDockerWorkerInterface) GetTrialObjValue(studyId string, tID string, objname string) (string, error) {
	cl, err := n.getConLog(tID, "")
	if err != nil {
		return "", err
	}
	for i := len(cl) - 1; i >= 0; i-- {
		ls := strings.Fields(cl[i])
		for _, l := range ls {
			v := strings.Split(l, "=")
			if v[0] == objname {
				return v[1], nil
			}
		}
	}
	return "", errors.New(fmt.Sprintf("No Objective Value Name %v  is found in log", objname))
}

func (n *NvDockerWorkerInterface) GetTrialEvLogs(studyId string, tID string, metrics []string, sinceTime string) ([]*api.EvaluationLog, error) {
	var ret []*api.EvaluationLog
	clog, err := n.getConLog(tID, sinceTime)
	if err != nil {
		return nil, err
	} else if clog == nil {
		return ret, nil
	}
	for _, ls := range clog {
		lsf := strings.Fields(ls)
		if len(lsf) == 0 {
			break
		}
		if lsf[0][8:] == sinceTime {
			continue
		}
		e := &api.EvaluationLog{Time: lsf[0][8:]}
		for _, l := range lsf[0:] {
			v := strings.Split(l, "=")
			for _, m := range metrics {
				if v[0] == m {
					e.Metrics = append(e.Metrics, &api.Metrics{Name: m, Value: v[1]})
				}
			}
		}
		if len(e.Metrics) > 0 {
			ret = append(ret, e)
		}
	}
	return ret, nil
}

func (n *NvDockerWorkerInterface) CheckRunningTrials(studyId string, objname string, metrics []string) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	if len(n.RunningTrialList[studyId]) == 0 {
		return nil
	}
	sc, _ := n.dbIf.GetStudyConfig(studyId)
	for _, t := range n.RunningTrialList[studyId] {
		status, err := n.dbIf.GetTrialStatus(t.TrialId)
		if err != nil {
			log.Printf("Error getting status of %s: %v", t.TrialId, err)
			continue
		}
		if status == api.TrialState_RUNNING {
			c, _ := n.IsTrialComplete(studyId, t.TrialId)
			var es []*api.EvaluationLog
			if len(t.EvalLogs) == 0 {
				es, err = n.GetTrialEvLogs(studyId, t.TrialId, metrics, "")
			} else {
				es, err = n.GetTrialEvLogs(studyId, t.TrialId, metrics, t.EvalLogs[len(t.EvalLogs)-1].Time)
			}
			if err != nil {
				log.Printf("GetTrialEvLogs Err %v", err)
				return err
			}
			if len(es) > 0 {
				t.EvalLogs = append(t.EvalLogs, es...)
			}
			if c {
				o, _ := n.GetTrialObjValue(studyId, t.TrialId, objname)
				t.ObjectiveValue = o
				t.Status = api.TrialState_COMPLETED
				mif := modeldb.ModelDbIF{}
				mr := &modeldb.ModelDbReq{
					Owner:          sc.Owner,
					Study:          sc.Name,
					Train:          t.TrialId,
					ModelPath:      "path/model",
					HyperParameter: make(map[string]string),
					Metrics:        make(map[string]float64),
				}
				of, err := strconv.ParseFloat(o, 64)
				if err != nil {
					log.Printf("ParseFloat err %v", err)
					return err
				}
				for _, p := range t.ParameterSet {
					mr.HyperParameter[strings.Replace(p.Name, "-", "", -1)] = p.Value
				}
				mr.Metrics[sc.ObjectiveValueName] = of
				for _, m := range sc.Metrics {
				MET_LABEL:
					for i := range t.EvalLogs {
						for _, em := range t.EvalLogs[len(t.EvalLogs)-1-i].Metrics {
							if em.Name == m {
								emv, err := strconv.ParseFloat(em.Value, 64)
								if err != nil {
									log.Printf("ParseFloat err %v", err)
									return err
								}
								mr.Metrics[m] = emv
								break MET_LABEL
							}
						}
					}
				}
				st, _ := time.Parse(time.RFC3339, t.EvalLogs[0].Time)
				et, _ := time.Parse(time.RFC3339, t.EvalLogs[len(t.EvalLogs)-1].Time)
				mr.Metrics["time-cost-Min"] = et.Sub(st).Minutes()
				mif.SendReq(mr)
				log.Printf("Trial %v is completed.", t.TrialId)
				log.Printf("Objective Value: %v", t.ObjectiveValue)
				n.CompletedTrialList[studyId] = append(n.CompletedTrialList[studyId], t)
				n.ngm.ReleaseGPU(t.TrialId)
				cid := n.tidToCid[t.TrialId]
				err = n.dcli.ContainerRemove(context.Background(), cid, types.ContainerRemoveOptions{Force: true})
				if err != nil {
					log.Printf("Container delete err %v", err)
				}
				delete(n.tidToCid, t.TrialId)
				if len(n.RunningTrialList[studyId]) <= 1 {
					n.RunningTrialList[studyId] = []*api.Trial{}
				} else {
					tn := t.TrialId
					for j, tt := range n.RunningTrialList[studyId] {
						if tt.TrialId == tn {
							n.RunningTrialList[studyId] = append(n.RunningTrialList[studyId][:j], n.RunningTrialList[studyId][j+1:]...)
							break
						}
					}
				}
			}
		}
	}
	return nil
}

func (n *NvDockerWorkerInterface) convertTrialToContainer(trials []*api.Trial, studyId string) []*container.Config {
	sc, _ := n.dbIf.GetStudyConfig(studyId)
	ret := make([]*container.Config, len(trials))
	n.mux.Lock()
	defer n.mux.Unlock()
	for i, t := range trials {
		n.PendingTrialList[studyId] = append(n.PendingTrialList[studyId], t)
		command := sc.Command
		for _, v := range t.ParameterSet {
			command = append(command, v.Name+"="+v.Value)
		}
		j := &container.Config{
			Image: sc.Image,
			Cmd:   command,
		}
		ret[i] = j
	}
	return ret
}

type trialQueueObj struct {
	StudyID string
	Trial   *api.Trial
	CConf   *container.Config
}

func (n *NvDockerWorkerInterface) schedulingLoop() error {
	tq := []*trialQueueObj{}
	t := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-n.stopSchedule:
			return nil
		case sr := <-n.addTrialQueCh:
			log.Printf("%v Trial set to trial queue", len(sr.Trials))
			for i, t := range sr.Trials {
				tq = append(tq, &trialQueueObj{
					StudyID: sr.StudyID,
					Trial:   t,
					CConf:   sr.CConf[i],
				})
			}
		case <-t.C:
			for i, t := range tq {
				sc, _ := n.dbIf.GetStudyConfig(t.StudyID)
				chc := &container.HostConfig{}
				if sc.Mount.Pvc != "" {
					chc.Binds = []string{sc.Mount.Pvc + ":" + sc.Mount.Path}
				}
				if sc.Gpu > 0 {
					ok, gid, err := n.ngm.AllocGPU(int(sc.Gpu), t.Trial.TrialId)
					if err != nil {
						log.Printf("AllocGPU error %v", err)
						break
					}
					if !ok {
						break
					}
					t.CConf.Env = []string{"NVIDIA_VISIBLE_DEVICES=" + strings.Join(gid, ",")}
					chc.Runtime = "nvidia"
				}
				resp, err := n.dcli.ContainerCreate(context.Background(), t.CConf, chc, nil, t.Trial.TrialId)
				if err != nil {
					log.Printf("Container create err %v", err)
					r, err := n.dcli.ImagePull(context.Background(), t.CConf.Image, types.ImagePullOptions{})
					if err != nil {
						log.Printf("Container create err %v", err)
						if sc.Gpu > 0 {
							n.ngm.ReleaseGPU(t.Trial.TrialId)
						}
						break
					}
					log.Printf("Image %v start to pull", t.CConf.Image)
					io.Copy(os.Stdout, r)
					resp, err = n.dcli.ContainerCreate(context.Background(), t.CConf, chc, nil, t.Trial.TrialId)
					if err != nil {
						log.Printf("Container create err %v", err)
						break
					}
				}
				if err := n.dcli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
					if sc.Gpu > 0 {
						n.ngm.ReleaseGPU(t.Trial.TrialId)
					}
					break
				}
				n.mux.Lock()
				n.tidToCid[t.Trial.TrialId] = resp.ID
				for i, pt := range n.PendingTrialList[t.StudyID] {
					if pt.TrialId == t.Trial.TrialId {
						if len(n.PendingTrialList[t.StudyID]) <= 1 {
							n.PendingTrialList[t.StudyID] = []*api.Trial{}
						} else {
							n.PendingTrialList[t.StudyID] = append(n.PendingTrialList[t.StudyID][:i], n.PendingTrialList[t.StudyID][i+1:]...)
						}
						n.dbIf.UpdateTrial(t.Trial.TrialId, api.TrialState_RUNNING)
						n.RunningTrialList[t.StudyID] = append(n.RunningTrialList[t.StudyID], pt)
						break
					}
				}
				n.mux.Unlock()
				if len(tq) <= 1 {
					tq = []*trialQueueObj{}
				} else {
					tq = append(tq[:i], tq[i+1:]...)
				}
				log.Printf("Trial %v spawned. container ID %v", t.Trial.TrialId, resp.ID)
			}
		case ds := <-n.deleteTrialQueCh:
			ntq := []*trialQueueObj{}
			for _, t := range tq {
				if ds != t.StudyID {
					ntq = append(ntq, t)
				}
			}
			tq = ntq
		}
	}
}

func (n *NvDockerWorkerInterface) SpawnWorkers(trials []*api.Trial, studyId string) error {
	n.addTrialQueCh <- spawnReq{StudyID: studyId, Trials: trials, CConf: n.convertTrialToContainer(trials, studyId)}
	return nil
}

func (n *NvDockerWorkerInterface) GetRunningTrials(studyId string) []*api.Trial {
	n.mux.Lock()
	defer n.mux.Unlock()
	return append(n.PendingTrialList[studyId], n.RunningTrialList[studyId]...)
}

func (n *NvDockerWorkerInterface) GetCompletedTrials(studyId string) []*api.Trial {
	n.mux.Lock()
	defer n.mux.Unlock()
	return n.CompletedTrialList[studyId]
}

func (n *NvDockerWorkerInterface) CleanWorkers(studyId string) error {
	n.mux.Lock()
	defer n.mux.Unlock()
	n.deleteTrialQueCh <- studyId
	for _, t := range n.RunningTrialList[studyId] {
		cid := n.tidToCid[t.TrialId]
		err := n.dcli.ContainerRemove(context.Background(), cid, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			log.Printf("Container delete err %v", err)
		}
		delete(n.tidToCid, t.TrialId)
		n.ngm.ReleaseGPU(t.TrialId)
	}
	delete(n.PendingTrialList, studyId)
	delete(n.RunningTrialList, studyId)
	return nil
}
