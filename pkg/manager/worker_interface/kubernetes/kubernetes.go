package kubernetes

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"

	"github.com/kubeflow/katib/api"
	"github.com/kubeflow/katib/pkg/db"
	"github.com/kubeflow/katib/pkg/earlystopping"
	wIF "github.com/kubeflow/katib/pkg/manager/worker_interface"
)

type KubernetesWorkerInterface struct {
	//Support MultiStudy
	RunningTrialList   map[string][]*api.Trial
	CompletedTrialList map[string][]*api.Trial
	clientset          *kubernetes.Clientset
	mux                *sync.Mutex
	db                 db.VizierDBInterface
}

func NewKubernetesWorkerInterface(cs *kubernetes.Clientset, db db.VizierDBInterface) *KubernetesWorkerInterface {
	return &KubernetesWorkerInterface{
		RunningTrialList:   make(map[string][]*api.Trial),
		CompletedTrialList: make(map[string][]*api.Trial),
		clientset:          cs,
		mux:                new(sync.Mutex),
		db:                 db,
	}
}

func (d *KubernetesWorkerInterface) convertTrialToManifest(trials []*api.Trial, tFile []byte, studyId string) []batchv1.Job {
	ret := make([]batchv1.Job, len(trials))
	BUFSIZE := 1024
	d.mux.Lock()
	defer d.mux.Unlock()
	for i, t := range trials {
		d.RunningTrialList[studyId] = append(d.RunningTrialList[studyId], t)
		k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(tFile), BUFSIZE).Decode(&ret[i])
		var args = []string{}
		for _, v := range t.ParameterSet {
			args = append(args, v.Name+"="+v.Value)
		}
		ret[i].ObjectMeta.Name = t.TrialId
		ret[i].Spec.Template.Spec.Containers[0].Name = t.TrialId + "-worker"
		ret[i].Spec.Template.Spec.Containers[0].Args = args
	}
	return ret
}

func (d *KubernetesWorkerInterface) storeTrialLog(tID string) error {
	pl, _ := d.clientset.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: "job-name=" + tID})
	if len(pl.Items) == 0 {
		return errors.New(fmt.Sprintf("No Pods are found in Job %v", tID))
	}

	mt, err := d.db.GetTrialTimestamp(tID)
	if err != nil {
		return err
	}
	logopt := apiv1.PodLogOptions{Timestamps: true}
	if mt != nil {
		logopt.SinceTime = &metav1.Time{Time: *mt}
	}

	logs, err := d.clientset.CoreV1().Pods(apiv1.NamespaceDefault).GetLogs(pl.Items[0].ObjectMeta.Name, &logopt).Do().Raw()
	if err != nil {
		return err
	}
	if len(logs) == 0 {
		return nil
	}
	err = d.db.StoreTrialLogs(tID, strings.Split(string(logs), "\n"))
	return err
}

func (d *KubernetesWorkerInterface) GetTrialObjValue(studyId string, tID string, objname string) (string, error) {
	return wIF.GetTrialObjValue(d.db, studyId, tID, objname)
}

func (d *KubernetesWorkerInterface) GetTrialEvLogs(studyId string, tID string, metrics []string, sinceTime string) ([]*api.EvaluationLog, error) {
	return wIF.GetTrialEvLogs(d.db, studyId, tID, metrics, sinceTime)
}

func (d *KubernetesWorkerInterface) PollingShouldStop(ess earlystopping.EarlyStoppingService, studyId string) chan bool {
	stop := make(chan bool)
	go func() {
		defer close(stop)
		tm := time.NewTimer(60 * time.Second)
		for {
			select {
			case <-tm.C:
				tm.Reset(60 * time.Second)
				//				d.mux.Lock()
				//				st := ess.ShouldTrialStop(d.RunningTrialList[studyId], d.CompletedTrialList[studyId], 10)
				//				jcl := d.clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
				//				pcl := d.clientset.CoreV1().Pods(apiv1.NamespaceDefault)
				//				for _, t := range st {
				//					jcl.Delete(t.TrialId, &metav1.DeleteOptions{})
				//					pl, _ := pcl.List(metav1.ListOptions{LabelSelector: "job-name=" + t.TrialId})
				//					pcl.Delete(pl.Items[0].ObjectMeta.Name, &metav1.DeleteOptions{})
				//					log.Printf("Trial %v is Killed.", t.TrialId)
				//					for i := range d.RunningTrialList[studyId] {
				//						if d.RunningTrialList[studyId][i].TrialId == t.TrialId {
				//							d.RunningTrialList[studyId][i].Status = api.TrialState_KILLED
				//							break
				//						}
				//					}
				//				}
				//				d.mux.Unlock()
			case <-stop:
				return
			}
		}
	}()
	return stop
}

func (d *KubernetesWorkerInterface) IsTrialComplete(studyId string, tID string) (bool, error) {
	jcl := d.clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
	ji, err := jcl.Get(tID, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if ji.Status.Succeeded == 0 {
		return false, nil
	}
	pl, _ := d.clientset.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: "job-name=" + tID})
	if len(pl.Items) == 0 {
		return false, errors.New(fmt.Sprintf("No Pods are found in Job %v", tID))
	}
	if pl.Items[0].Status.Phase == "Succeeded" {
		return true, nil
	}
	return false, nil
}

func (d *KubernetesWorkerInterface) CheckRunningTrials(studyId string, objname string) error {
	allcomp := true
	d.mux.Lock()
	sc, _ := d.db.GetStudyConfig(studyId)
	metrics := sc.Metrics
	defer d.mux.Unlock()
	if len(d.RunningTrialList[studyId]) == 0 {
		return nil
	}
	for i, t := range d.RunningTrialList[studyId] {
		status, err := d.db.GetTrialStatus(t.TrialId)
		if err != nil {
			log.Printf("Error getting status of %s: %v", t.TrialId, err)
			continue
		}
		if status == api.TrialState_RUNNING {
			c, err := d.IsTrialComplete(studyId, t.TrialId)
			if err != nil {
				log.Printf("IsTrialComplete: %v", err)
			}
			err = d.storeTrialLog(t.TrialId)
			if err != nil {
				log.Printf("Error storing trial log of %s: %v", t.TrialId, err)
			}
			// TODO: remove the following code block after
			// suggestions and earlystopping switch to
			// read those values from DB
			if c {
				o, _ := d.GetTrialObjValue(studyId, t.TrialId, objname)
				d.RunningTrialList[studyId][i].ObjectiveValue = o
				d.RunningTrialList[studyId][i].Status = api.TrialState_COMPLETED
			} else {
				allcomp = false
				var es []*api.EvaluationLog
				if len(d.RunningTrialList[studyId][i].EvalLogs) == 0 {
					es, _ = d.GetTrialEvLogs(studyId, t.TrialId, metrics, "")
				} else {
					es, _ = d.GetTrialEvLogs(studyId, t.TrialId, metrics, d.RunningTrialList[studyId][i].EvalLogs[len(d.RunningTrialList[studyId][i].EvalLogs)-1].Time)
				}
				if len(es) > 0 {
					d.RunningTrialList[studyId][i].EvalLogs = append(d.RunningTrialList[studyId][i].EvalLogs, es...)
				}
			}
		} else if status == api.TrialState_PENDING {
			allcomp = false
		}
	}
	if allcomp {
		for i, t := range d.RunningTrialList[studyId] {
			log.Printf("%v is completed.", t.TrialId)
			log.Printf("Objective Value: %v", d.RunningTrialList[studyId][i].ObjectiveValue)
			log.Printf("Tags: %v", t.Tags)
			//			for _, l := range t.EvalLogs {
			//				log.Printf("\tEval Logs %v %v\n", l.Time, l.Value)
			//			}
		}
		d.CompletedTrialList[studyId] = append(d.CompletedTrialList[studyId], d.RunningTrialList[studyId]...)
		d.RunningTrialList[studyId] = []*api.Trial{}
		return nil
	}
	return nil
}

func (d *KubernetesWorkerInterface) SpawnWorkers(trials []*api.Trial, studyId string) error {
	// Notice: Missing in the repo.
	tFile, _ := ioutil.ReadFile("/conf/template.yml")
	jobs := d.convertTrialToManifest(trials, tFile, studyId)
	jcl := d.clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
	for _, j := range jobs {
		result, err := jcl.Create(&j)
		if err != nil {
			return err
		}
		err = d.db.UpdateTrial(j.ObjectMeta.Name, api.TrialState_RUNNING)
		if err != nil {
			log.Printf("Error updating status for %s: %v", j.ObjectMeta.Name, err)
		}

		log.Printf("Created Job %q.", result.GetObjectMeta().GetName())
	}
	return nil
}

func (d *KubernetesWorkerInterface) GetRunningTrials(studyId string) []*api.Trial {
	return d.RunningTrialList[studyId]
}

func (d *KubernetesWorkerInterface) GetCompletedTrials(studyId string) []*api.Trial {
	return d.CompletedTrialList[studyId]
}

func (d *KubernetesWorkerInterface) CleanWorkers(studyId string) error {
	jcl := d.clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
	pcl := d.clientset.CoreV1().Pods(apiv1.NamespaceDefault)
	for _, t := range d.RunningTrialList[studyId] {
		jcl.Delete(t.TrialId, &metav1.DeleteOptions{})
		pl, _ := pcl.List(metav1.ListOptions{LabelSelector: "job-name=" + t.TrialId})
		pcl.Delete(pl.Items[0].ObjectMeta.Name, &metav1.DeleteOptions{})
	}
	for _, t := range d.CompletedTrialList[studyId] {
		jcl.Delete(t.TrialId, &metav1.DeleteOptions{})
		pl, _ := pcl.List(metav1.ListOptions{LabelSelector: "job-name=" + t.TrialId})
		pcl.Delete(pl.Items[0].ObjectMeta.Name, &metav1.DeleteOptions{})
	}
	delete(d.RunningTrialList, studyId)
	delete(d.CompletedTrialList, studyId)
	return nil
}
func (d *KubernetesWorkerInterface) CompleteTrial(studyId string, tID string, iscomplete bool) error {
	jcl := d.clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
	pcl := d.clientset.CoreV1().Pods(apiv1.NamespaceDefault)
	jcl.Delete(tID, &metav1.DeleteOptions{})
	pl, _ := pcl.List(metav1.ListOptions{LabelSelector: "job-name=" + tID})
	pcl.Delete(pl.Items[0].ObjectMeta.Name, &metav1.DeleteOptions{})
	return nil
}
