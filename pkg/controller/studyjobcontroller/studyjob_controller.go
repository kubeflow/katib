/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package studyjobcontroller

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"text/template"

	"github.com/kubeflow/katib/pkg"
	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new StudyJobController Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this studyjobcontroller.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStudyJobController{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("studyjob-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to StudyJobController
	err = c.Watch(&source.Kind{Type: &katibv1alpha1.StudyJob{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by StudyJobController - change this for objects you create
	err = c.Watch(&source.Kind{Type: &katibv1alpha1.StudyJob{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &katibv1alpha1.StudyJob{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(
		&source.Kind{Type: &batchv1.Job{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &katibv1alpha1.StudyJob{},
		})
	if err != nil {
		return err
	}

	err = c.Watch(
		&source.Kind{Type: &batchv1beta.CronJob{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &katibv1alpha1.StudyJob{},
		})
	if err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileStudyJobController{}

// ReconcileStudyJobController reconciles a StudyJob object
type ReconcileStudyJobController struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a StudyJob object and makes changes based on the state read
// and what is in the StudyJob.Spec
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=studyjob.kubeflow.org,resources=studyjob,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileStudyJobController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the StudyJob instance
	instance := &katibv1alpha1.StudyJob{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			log.Println("No instance")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	switch instance.Status.Condition {
	case katibv1alpha1.ConditionCompleted:
		err = r.checkStatus(instance, request.Namespace)
	case katibv1alpha1.ConditionFailed:
		err = r.checkStatus(instance, request.Namespace)
	case katibv1alpha1.ConditionRunning:
		err = r.checkStatus(instance, request.Namespace)
	default:
		err = r.initializeStudy(instance, request.Namespace)
		if err != nil {
			r.Update(context.TODO(), instance)
			log.Printf("Fail to initialize %v", err)
			return reconcile.Result{}, err
		}
	}
	if err != nil {
		r.Update(context.TODO(), instance)
		log.Printf("Fail to check status %v", err)
		return reconcile.Result{}, err
	}
	err = r.Update(context.TODO(), instance)
	return reconcile.Result{}, err
}

func (r *ReconcileStudyJobController) getStudyConf(instance *katibv1alpha1.StudyJob) (*katibapi.StudyConfig, error) {
	sconf := &katibapi.StudyConfig{
		Metrics: []string{},
		ParameterConfigs: &katibapi.StudyConfig_ParameterConfigs{
			Configs: []*katibapi.ParameterConfig{},
		},
	}
	sconf.Name = instance.Spec.StudyName
	sconf.Owner = instance.Spec.Owner
	if instance.Spec.OptimizationGoal != nil {

		sconf.OptimizationGoal = *instance.Spec.OptimizationGoal
	}
	sconf.ObjectiveValueName = instance.Spec.ObjectiveValueName
	switch instance.Spec.OptimizationType {
	case katibv1alpha1.OptimizationTypeMinimize:
		sconf.OptimizationType = katibapi.OptimizationType_MINIMIZE
	case katibv1alpha1.OptimizationTypeMaximize:
		sconf.OptimizationType = katibapi.OptimizationType_MAXIMIZE
	default:
		sconf.OptimizationType = katibapi.OptimizationType_UNKNOWN_OPTIMIZATION
	}
	for _, m := range instance.Spec.MetricsNames {
		sconf.Metrics = append(sconf.Metrics, m)
	}
	for _, pc := range instance.Spec.ParameterConfigs {
		p := &katibapi.ParameterConfig{
			Feasible: &katibapi.FeasibleSpace{},
		}
		p.Name = pc.Name
		p.Feasible.Min = pc.Feasible.Min
		p.Feasible.Max = pc.Feasible.Max
		p.Feasible.List = pc.Feasible.List
		switch pc.ParameterType {
		case katibv1alpha1.ParameterTypeUnknown:
			p.ParameterType = katibapi.ParameterType_UNKNOWN_TYPE
		case katibv1alpha1.ParameterTypeDouble:
			p.ParameterType = katibapi.ParameterType_DOUBLE
		case katibv1alpha1.ParameterTypeInt:
			p.ParameterType = katibapi.ParameterType_INT
		case katibv1alpha1.ParameterTypeDiscrete:
			p.ParameterType = katibapi.ParameterType_DISCRETE
		case katibv1alpha1.ParameterTypeCategorical:
			p.ParameterType = katibapi.ParameterType_CATEGORICAL
		}
		sconf.ParameterConfigs.Configs = append(sconf.ParameterConfigs.Configs, p)
	}
	sconf.JobId = string(instance.UID)
	return sconf, nil
}

func (r *ReconcileStudyJobController) checkGoal(instance *katibv1alpha1.StudyJob, c katibapi.ManagerClient, wids []string) (bool, error) {
	if instance.Spec.OptimizationGoal == nil {
		return false, nil
	}
	getMetricsRequest := &katibapi.GetMetricsRequest{
		StudyId:   instance.Status.StudyId,
		WorkerIds: wids,
	}
	mr, err := c.GetMetrics(context.Background(), getMetricsRequest)
	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == codes.ResourceExhausted {
			//GRPC Message is probably too long, get metrics one by one
			var mls []*katibapi.MetricsLogSet
			for _, wid := range wids {
				getMetricsRequest := &katibapi.GetMetricsRequest{
					StudyId:   instance.Status.StudyId,
					WorkerIds: []string{wid},
				}
				cmr, err := c.GetMetrics(context.Background(), getMetricsRequest)
				if err != nil {
					return false, err
				}
				mls = append(mls, cmr.MetricsLogSets...)
			}
			mr.MetricsLogSets = mls
		} else {
			//Unknown error
			return false, err
		}
	}
	goal := false
	for _, mls := range mr.MetricsLogSets {
		for _, ml := range mls.MetricsLogs {
			if ml.Name == instance.Spec.ObjectiveValueName {
				if len(ml.Values) > 0 {
					curValue, _ := strconv.ParseFloat(ml.Values[len(ml.Values)-1].Value, 32)
					if instance.Spec.OptimizationType == katibv1alpha1.OptimizationTypeMinimize {
						if curValue < *instance.Spec.OptimizationGoal {
							goal = true
						}
						if instance.Status.BestObjctiveValue != nil {
							if *instance.Status.BestObjctiveValue > curValue {
								instance.Status.BestObjctiveValue = &curValue
							}
						} else {
							instance.Status.BestObjctiveValue = &curValue
						}
						for i := range instance.Status.Trials {
							for j := range instance.Status.Trials[i].WorkerList {
								if instance.Status.Trials[i].WorkerList[j].WorkerId == mls.WorkerId {
									instance.Status.Trials[i].WorkerList[j].ObjctiveValue = &curValue
								}
							}
						}
					} else if instance.Spec.OptimizationType == katibv1alpha1.OptimizationTypeMaximize {
						if curValue > *instance.Spec.OptimizationGoal {
							goal = true
						}
						if instance.Status.BestObjctiveValue != nil {
							if *instance.Status.BestObjctiveValue < curValue {
								instance.Status.BestObjctiveValue = &curValue
							}
						} else {
							instance.Status.BestObjctiveValue = &curValue
						}
						for i := range instance.Status.Trials {
							for j := range instance.Status.Trials[i].WorkerList {
								if instance.Status.Trials[i].WorkerList[j].WorkerId == mls.WorkerId {
									instance.Status.Trials[i].WorkerList[j].ObjctiveValue = &curValue
								}
							}
						}
					}
				}
				break
			}
		}
	}
	return goal, nil
}

func (r *ReconcileStudyJobController) initializeStudy(instance *katibv1alpha1.StudyJob, ns string) error {
	if instance.Spec.SuggestionSpec == nil {
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		return nil
	}
	if instance.Spec.SuggestionSpec.SuggestionAlgorithm == "" {
		instance.Spec.SuggestionSpec.SuggestionAlgorithm = "random"
	}
	instance.Status.Condition = katibv1alpha1.ConditionRunning

	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect katib manager error %v", err)
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		return nil
	}
	defer conn.Close()
	c := katibapi.NewManagerClient(conn)

	studyConfig, err := r.getStudyConf(instance)
	if err != nil {
		return err
	}

	log.Printf("Create Study %s", studyConfig.Name)
	//CreateStudy
	studyId, err := r.createStudy(c, studyConfig)
	if err != nil {
		return err
	}
	instance.Status.StudyId = studyId
	log.Printf("Study: %s Suggestion Spec %v", studyId, instance.Spec.SuggestionSpec)
	var sspec *katibv1alpha1.SuggestionSpec
	if instance.Spec.SuggestionSpec != nil {
		sspec = instance.Spec.SuggestionSpec
	} else {
		sspec = &katibv1alpha1.SuggestionSpec{}
	}
	sspec.SuggestionParameters = append(sspec.SuggestionParameters,
		katibapi.SuggestionParameter{
			Name:  "SuggestionCount",
			Value: "0",
		})
	sPID, err := r.setSuggestionParam(c, studyId, sspec)
	if err != nil {
		return err
	}
	instance.Status.SuggestionParameterId = sPID
	instance.Status.SuggestionCount += 1
	instance.Status.Condition = katibv1alpha1.ConditionRunning
	return nil
}

func (r *ReconcileStudyJobController) checkStatus(instance *katibv1alpha1.StudyJob, ns string) error {
	nextSuggestionSchedule := true
	var cwids []string
	if instance.Status.Condition == katibv1alpha1.ConditionCompleted || instance.Status.Condition == katibv1alpha1.ConditionFailed {
		nextSuggestionSchedule = false
	}
	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect katib manager error %v", err)
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		return nil
	}
	defer conn.Close()
	c := katibapi.NewManagerClient(conn)
	for i, t := range instance.Status.Trials {
		for j, w := range t.WorkerList {
			if w.Condition == katibv1alpha1.ConditionCompleted {
				if w.ObjctiveValue == nil {
					cwids = append(cwids, w.WorkerId)
				}
				continue
			}
			nextSuggestionSchedule = false
			switch w.Kind {
			case "Job":
				job := &batchv1.Job{}
				nname := types.NamespacedName{Namespace: ns, Name: w.WorkerId}
				joberr := r.Client.Get(context.TODO(), nname, job)
				if joberr != nil {
					continue
				}
				cjob := &batchv1beta.CronJob{}
				cjoberr := r.Client.Get(context.TODO(), nname, cjob)
				if job.Status.Active == 0 && job.Status.Succeeded > 0 {
					ctime := job.Status.CompletionTime
					if cjoberr == nil {
						if ctime != nil && cjob.Status.LastScheduleTime != nil {
							if ctime.Before(cjob.Status.LastScheduleTime) {
								instance.Status.Trials[i].WorkerList[j].Condition = katibv1alpha1.ConditionCompleted
								susp := true
								cjob.Spec.Suspend = &susp
								if err := r.Update(context.TODO(), cjob); err != nil {
									return err
								}
								cwids = append(cwids, w.WorkerId)
							}
						}
					}
				} else if job.Status.Active > 0 {
					instance.Status.Trials[i].WorkerList[j].Condition = katibv1alpha1.ConditionRunning
					if errors.IsNotFound(cjoberr) {
						r.spawnMetricsCollector(instance, c, instance.Status.StudyId, t.TrialId, w.WorkerId, ns, instance.Spec.MetricsCollectorSpec)
					}
				} else if job.Status.Failed > 0 {
					instance.Status.Trials[i].WorkerList[j].Condition = katibv1alpha1.ConditionFailed
				}
			}
		}
	}
	if len(cwids) > 0 {
		goal, err := r.checkGoal(instance, c, cwids)
		if goal {
			log.Printf("Study %s reached to the goal. It is completed", instance.Status.StudyId)
			instance.Status.Condition = katibv1alpha1.ConditionCompleted
			nextSuggestionSchedule = false
		}
		if err != nil {
			log.Printf("Check Goal failed %v", err)
		}
	}
	if nextSuggestionSchedule {
		return r.getAndRunSuggestion(instance, c, ns)
	} else {
		return nil
	}
}

func (r *ReconcileStudyJobController) getAndRunSuggestion(instance *katibv1alpha1.StudyJob, c katibapi.ManagerClient, ns string) error {
	//Check Suggestion Count
	sps, err := r.getSuggestionParam(c, instance.Status.SuggestionParameterId)
	if err != nil {
		return err
	}
	for i := range sps {
		if sps[i].Name == "SuggestionCount" {
			count, _ := strconv.Atoi(sps[i].Value)
			if count >= instance.Status.SuggestionCount+1 {
				return fmt.Errorf("Suggestion count mismatched. May be duplicate suggestion request")
			}
			sps[i].Value = strconv.Itoa(instance.Status.SuggestionCount + 1)
		}
	}
	//GetSuggestion
	getSuggestReply, err := r.getSuggestion(
		c,
		instance.Status.StudyId,
		instance.Spec.SuggestionSpec,
		instance.Status.SuggestionParameterId)
	if err != nil {
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		return err
	}
	trials := getSuggestReply.Trials
	if len(trials) <= 0 {
		log.Printf("Study %s is completed", instance.Status.StudyId)
		instance.Status.Condition = katibv1alpha1.ConditionCompleted
		return nil
	}
	log.Printf("Study: %s Suggestions %v", instance.Status.StudyId, getSuggestReply)
	wkind, err := r.getWorkerKind(instance.Spec.WorkerSpec)
	if err != nil {
		log.Printf("getWorkerKind error %v", err)
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		return err
	}
	for _, t := range trials {
		wid, err := r.spawnWorker(instance, c, instance.Status.StudyId, t, instance.Spec.WorkerSpec, wkind, false)
		if err != nil {
			log.Printf("Spawn worker error %v", err)
			instance.Status.Condition = katibv1alpha1.ConditionFailed
			return err
		}
		instance.Status.Trials = append(
			instance.Status.Trials,
			katibv1alpha1.TrialSet{
				TrialId: t.TrialId,
				WorkerList: []katibv1alpha1.WorkerCondition{
					katibv1alpha1.WorkerCondition{
						WorkerId:  wid,
						Kind:      wkind,
						Condition: katibv1alpha1.ConditionCreated,
					},
				},
			},
		)
	}
	//Update Suggestion Count
	sspr := &katibapi.SetSuggestionParametersRequest{
		StudyId:              instance.Status.StudyId,
		SuggestionAlgorithm:  instance.Spec.SuggestionSpec.SuggestionAlgorithm,
		ParamId:              instance.Status.SuggestionParameterId,
		SuggestionParameters: sps,
	}
	_, err = c.SetSuggestionParameters(context.Background(), sspr)
	if err != nil {
		log.Printf("Study %s Suggestion Count update Error %v", instance.Status.StudyId, err)
		return err
	}
	instance.Status.SuggestionCount += 1
	return nil
}

func (r *ReconcileStudyJobController) createStudy(c katibapi.ManagerClient, studyConfig *katibapi.StudyConfig) (string, error) {
	ctx := context.Background()
	createStudyreq := &katibapi.CreateStudyRequest{
		StudyConfig: studyConfig,
	}
	createStudyreply, err := c.CreateStudy(ctx, createStudyreq)
	if err != nil {
		log.Printf("CreateStudy Error %v", err)
		return "", err
	}
	studyId := createStudyreply.StudyId
	log.Printf("Study ID %s", studyId)
	getStudyreq := &katibapi.GetStudyRequest{
		StudyId: studyId,
	}
	getStudyReply, err := c.GetStudy(ctx, getStudyreq)
	if err != nil {
		log.Printf("Study: %s GetConfig Error %v", studyId, err)
		return "", err
	}
	log.Printf("Study ID %s StudyConf%v", studyId, getStudyReply.StudyConfig)
	return studyId, nil
}

func (r *ReconcileStudyJobController) setSuggestionParam(c katibapi.ManagerClient, studyId string, suggestionSpec *katibv1alpha1.SuggestionSpec) (string, error) {
	ctx := context.Background()
	pid := ""
	if suggestionSpec.SuggestionParameters != nil {
		sspr := &katibapi.SetSuggestionParametersRequest{
			StudyId:             studyId,
			SuggestionAlgorithm: suggestionSpec.SuggestionAlgorithm,
		}
		for _, p := range suggestionSpec.SuggestionParameters {
			sspr.SuggestionParameters = append(
				sspr.SuggestionParameters,
				&katibapi.SuggestionParameter{
					Name:  p.Name,
					Value: p.Value,
				},
			)
		}
		setSuggesitonParameterReply, err := c.SetSuggestionParameters(ctx, sspr)
		if err != nil {
			log.Printf("Study %s SetConfig Error %v", studyId, err)
			return "", err
		}
		log.Printf("Study: %s setSuggesitonParameterReply %v", studyId, setSuggesitonParameterReply)
		pid = setSuggesitonParameterReply.ParamId
	}
	return pid, nil
}

func (r *ReconcileStudyJobController) getSuggestionParam(c katibapi.ManagerClient, paramId string) ([]*katibapi.SuggestionParameter, error) {
	ctx := context.Background()
	gsreq := &katibapi.GetSuggestionParametersRequest{
		ParamId: paramId,
	}
	gsrep, err := c.GetSuggestionParameters(ctx, gsreq)
	return gsrep.SuggestionParameters, err
}
func (r *ReconcileStudyJobController) getSuggestion(c katibapi.ManagerClient, studyId string, suggestionSpec *katibv1alpha1.SuggestionSpec, sParamID string) (*katibapi.GetSuggestionsReply, error) {
	ctx := context.Background()
	getSuggestRequest := &katibapi.GetSuggestionsRequest{
		StudyId:             studyId,
		SuggestionAlgorithm: suggestionSpec.SuggestionAlgorithm,
		RequestNumber:       int32(suggestionSpec.RequestNumber),
		//RequestNumber=0 means get all grids.
		ParamId: sParamID,
	}
	getSuggestReply, err := c.GetSuggestions(ctx, getSuggestRequest)
	if err != nil {
		log.Printf("Study: %s GetSuggestion Error %v", studyId, err)
		return nil, err
	}
	log.Printf("Study: %s CreatedTrials :", studyId)
	for _, t := range getSuggestReply.Trials {
		log.Printf("\t%v", t)
	}
	return getSuggestReply, nil
}

type WorkerInstance struct {
	StudyId         string
	TrialId         string
	WorkerId        string
	HyperParameters []*katibapi.Parameter
}

func (r *ReconcileStudyJobController) getWorkerKind(workerSpec *katibv1alpha1.WorkerSpec) (string, error) {
	var typeChecker interface{}
	BUFSIZE := 1024
	_, m, err := r.getWorkerManifest(
		nil,
		"validation",
		&katibapi.Trial{
			TrialId:      "validation",
			ParameterSet: []*katibapi.Parameter{},
		},
		workerSpec,
		"",
		true,
	)
	if err != nil {
		return "", err
	}
	if err := k8syaml.NewYAMLOrJSONDecoder(m, BUFSIZE).Decode(&typeChecker); err != nil {
		log.Printf("Yaml decode validation error %v", err)
		return "", err
	}
	tcMap, ok := typeChecker.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	wkind, ok := tcMap["kind"]
	if !ok {
		return "", fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	wkind_s, ok := wkind.(string)
	if !ok {
		return "", fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	return wkind_s, nil
}

func (r *ReconcileStudyJobController) spawnWorker(instance *katibv1alpha1.StudyJob, c katibapi.ManagerClient, studyId string, trial *katibapi.Trial, workerSpec *katibv1alpha1.WorkerSpec, wkind string, dryrun bool) (string, error) {
	wid, wm, err := r.getWorkerManifest(c, studyId, trial, workerSpec, wkind, false)
	if err != nil {
		return "", err
	}
	BUFSIZE := 1024
	switch wkind {
	case "Job":
		var job batchv1.Job
		if err := k8syaml.NewYAMLOrJSONDecoder(wm, BUFSIZE).Decode(&job); err != nil {
			instance.Status.Condition = katibv1alpha1.ConditionFailed
			log.Printf("Yaml decode error %v", err)
			return "", err
		}
		if err := controllerutil.SetControllerReference(instance, &job, r.scheme); err != nil {
			instance.Status.Condition = katibv1alpha1.ConditionFailed
			log.Printf("SetControllerReference error %v", err)
			return "", err
		}
		if err := r.Create(context.TODO(), &job); err != nil {
			instance.Status.Condition = katibv1alpha1.ConditionFailed
			log.Printf("Job Create error %v", err)
			return "", err
		}
	}
	return wid, nil
}

func (r *ReconcileStudyJobController) getWorkerManifest(c katibapi.ManagerClient, studyId string, trial *katibapi.Trial, workerSpec *katibv1alpha1.WorkerSpec, kind string, dryrun bool) (string, *bytes.Buffer, error) {
	var wtp *template.Template
	var err error
	wpath := "/worker-template/defaultWorkerTemplate.yaml"

	if workerSpec != nil {
		if workerSpec.GoTemplate.Path != "" {
			wpath = workerSpec.GoTemplate.Path
		}
	}
	wtp, err = template.ParseFiles(wpath)
	if err != nil {
		return "", nil, err
	}
	var wid string
	if dryrun {
		wid = "validation"
	} else {
		cwreq := &katibapi.RegisterWorkerRequest{
			Worker: &katibapi.Worker{
				StudyId: studyId,
				TrialId: trial.TrialId,
				Status:  katibapi.State_PENDING,
				Type:    kind,
			},
		}
		cwrep, err := c.RegisterWorker(context.Background(), cwreq)
		if err != nil {
			return "", nil, err
		}
		wid = cwrep.WorkerId
	}

	wi := WorkerInstance{
		StudyId:  studyId,
		TrialId:  trial.TrialId,
		WorkerId: wid,
	}
	var b bytes.Buffer
	for _, p := range trial.ParameterSet {
		wi.HyperParameters = append(wi.HyperParameters, p)
	}
	err = wtp.Execute(&b, wi)
	if err != nil {
		return "", nil, err
	}
	return wid, &b, nil
}

func (r *ReconcileStudyJobController) spawnMetricsCollector(instance *katibv1alpha1.StudyJob, c katibapi.ManagerClient, studyId string, trialId string, workerId string, namespace string, mcs *katibv1alpha1.MetricsCollectorSpec) error {
	var mcjob batchv1beta.CronJob
	BUFSIZE := 1024
	mcm, err := r.getMetricsCollectorManifest(studyId, trialId, workerId, namespace, mcs)
	if err != nil {
		log.Printf("getMetricsCollectorManifest error %v", err)
		return err
	}

	if err := k8syaml.NewYAMLOrJSONDecoder(mcm, BUFSIZE).Decode(&mcjob); err != nil {
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		log.Printf("MetricsCollector Yaml decode error %v", err)
		return err
	}

	if err := controllerutil.SetControllerReference(instance, &mcjob, r.scheme); err != nil {
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		log.Printf("MetricsCollector SetControllerReference error %v", err)
		return err
	}

	if err := r.Create(context.TODO(), &mcjob); err != nil {
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		log.Printf("MetricsCollector Job Create error %v", err)
		return err
	}
	return nil
}

func (r *ReconcileStudyJobController) getMetricsCollectorManifest(studyId string, trialId string, workerId string, namespace string, mcs *katibv1alpha1.MetricsCollectorSpec) (*bytes.Buffer, error) {
	var mtp *template.Template
	var err error
	tmpValues := map[string]string{
		"StudyId":   studyId,
		"TrialId":   trialId,
		"WorkerId":  workerId,
		"NameSpace": namespace,
	}
	mpath := "/metricscollector-template/defaultMetricsCollectorTemplate.yaml"
	if mcs != nil {
		if mcs.GoTemplate.Path != "" {
			mpath = mcs.GoTemplate.Path
		}
	}
	mtp, err = template.ParseFiles(mpath)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	err = mtp.Execute(&b, tmpValues)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
