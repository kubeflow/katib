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
	"math/rand"
	"strconv"
	"text/template"

	"github.com/kubeflow/katib/pkg"
	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
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
	log.Println("c.Watch(&source.Kind{Type: &katibv1alpha1.StudyJobController{}}, &handler.EnqueueRequestForOwner{")
	err = c.Watch(&source.Kind{Type: &katibv1alpha1.StudyJob{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &katibv1alpha1.StudyJob{},
	})
	if err != nil {
		return err
	}
	log.Println("add comp")

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
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	if instance.Status.Condition == katibv1alpha1.ConditionRunning || instance.Status.Condition == katibv1alpha1.ConditionCompleted || instance.Status.Condition == katibv1alpha1.ConditionFailed {
		return reconcile.Result{}, nil
	}
	r.controllerloop(instance)
	return reconcile.Result{}, nil
}

func (r *ReconcileStudyJobController) getStudyConf(instance *katibv1alpha1.StudyJob) (*katibapi.StudyConfig, error) {

	sconf := &katibapi.StudyConfig{
		Metrics: []string{},
		ParameterConfigs: &katibapi.StudyConfig_ParameterConfigs{
			Configs: []*katibapi.ParameterConfig{},
		},
	}
	if instance.Spec.StudySpec.Name == "" {
		return nil, fmt.Errorf("StudyName must be set")
	}

	sconf.Name = instance.Spec.StudySpec.Name
	sconf.Owner = instance.Spec.StudySpec.Owner
	if instance.Spec.StudySpec.OptimizationGoal != nil {

		sconf.OptimizationGoal = *instance.Spec.StudySpec.OptimizationGoal
	}
	sconf.ObjectiveValueName = instance.Spec.StudySpec.ObjectiveValueName
	switch instance.Spec.StudySpec.OptimizationType {
	case katibv1alpha1.OptimizationTypeMinimize:
		sconf.OptimizationType = katibapi.OptimizationType_MINIMIZE
	case katibv1alpha1.OptimizationTypeMaximize:
		sconf.OptimizationType = katibapi.OptimizationType_MAXIMIZE
	default:
		sconf.OptimizationType = katibapi.OptimizationType_UNKNOWN_OPTIMIZATION
	}
	for _, m := range instance.Spec.StudySpec.MetricsNames {
		sconf.Metrics = append(sconf.Metrics, m)
	}
	for _, pc := range instance.Spec.StudySpec.ParameterConfigs {
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
	return sconf, nil
}

func (r *ReconcileStudyJobController) checkGoal(instance *katibv1alpha1.StudyJob, sconf *katibapi.StudyConfig, mr *katibapi.GetMetricsReply) bool {
	if instance.Spec.StudySpec.OptimizationGoal == nil {
		return false
	}
	for _, mls := range mr.MetricsLogSets {
		for _, ml := range mls.MetricsLogs {
			if ml.Name == sconf.ObjectiveValueName {
				curValue, _ := strconv.ParseFloat(ml.Values[len(ml.Values)-1].Value, 32)
				if sconf.OptimizationType == katibapi.OptimizationType_MINIMIZE && curValue < sconf.OptimizationGoal {
					return true
				} else if sconf.OptimizationType == katibapi.OptimizationType_MAXIMIZE && curValue > sconf.OptimizationGoal {
					return true
				} else {
					return false
				}
			}
		}
	}
	return false
}

//Main loop of StudyJob Controller.
//This loop is for each StudyJob Object.
//Create Study, set Suggesiton Parameters, GetSuggestion and RunTrials. Then polling the workers.
func (r *ReconcileStudyJobController) controllerloop(instance *katibv1alpha1.StudyJob) {
	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect katib manager error %v", err)
		return
	}
	defer conn.Close()
	if instance.Spec.StudySpec == nil || instance.Spec.SuggestionSpec == nil {
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		r.Update(context.TODO(), instance)
		return
	}
	if instance.Spec.SuggestionSpec.SuggestionAlgorithm == "" {
		instance.Spec.SuggestionSpec.SuggestionAlgorithm = "random"
	}
	instance.Status.Condition = katibv1alpha1.ConditionRunning
	if err := r.Update(context.TODO(), instance); err != nil {
		return
	}
	//ctx := context.Background()
	c := katibapi.NewManagerClient(conn)

	studyConfig, err := r.getStudyConf(instance)
	if err != nil {
		return
	}

	log.Printf("Create Study %s", studyConfig.Name)
	//CreateStudy
	studyId, err := r.createStudy(c, studyConfig)
	if err != nil {
		return
	}
	instance.Status.StudyId = studyId
	if err := r.Update(context.TODO(), instance); err != nil {
		return
	}

	log.Printf("Study: %s Get Suggestion %v", studyId, instance.Spec.SuggestionSpec)
	sPID, err := r.setSuggestionParam(c, studyId, instance.Spec.SuggestionSpec)
	if err != nil {
		return
	}
	//	for true {
	//GetSuggestion
	getSuggestReply, err := r.getSuggestion(c, studyId, instance.Spec.SuggestionSpec, sPID)
	if err != nil {
		return
	}
	trials := getSuggestReply.Trials
	if len(trials) <= 0 {
		log.Printf("Study %s is completed", studyId)
		return
	}
	log.Printf("Study: %s Suggestions %v", studyId, getSuggestReply)
	wids, wins, err := r.getWorerManifest(instance.Spec.WorkerSpec, trials)
	if err != nil {
		log.Printf("getWorerManifest error %v", err)
		return
	}
	jobs := make([]batchv1.Job, len(wins))
	BUFSIZE := 1024
	for i := range jobs {
		log.Printf("Manifest %s", wins[i].String())
		err = k8syaml.NewYAMLOrJSONDecoder(wins[i], BUFSIZE).Decode(&jobs[i])
		if err != nil {
			log.Printf("Yaml decode error %v", err)
			log.Printf("Manifest %s", wins[i].String())
			return
		}
		log.Printf("WorkerID: %s\n Manifest %v\n", wids[i], jobs[i])
	}
	//workerIds, tsl, err := r.runTrial(c, studyId, trials, studyConfig, instance.Spec.WorkerSpec)
	//if err != nil {
	//		return
	//	}
	//instance.Status.Trials = append(instance.Status.Trials, trials...)
	if err := r.Update(context.TODO(), instance); err != nil {
		return
	}
	//		getMetricsReply := &katibapi.GetMetricsReply{}
	//		for true {
	//			time.Sleep(10 * time.Second)
	//			getMetricsRequest := &katibapi.GetMetricsRequest{
	//				StudyId:   studyId,
	//				WorkerIds: workerIds,
	//			}
	//			//GetMetrics
	//			getMetricsReply, err = c.GetMetrics(ctx, getMetricsRequest)
	//			if err != nil {
	//				continue
	//			}
	//			//Save or Update model on ModelDB
	//			r.saveOrUpdateModel(c, getMetricsReply, studyConfig, studyId)
	//			if r.isCompletedAllWorker(c, getMetricsReply.MetricsLogSets) {
	//				break
	//			}
	//		}
	//		if instance.Spec.SuggestionSpec.SuggestionAlgorithm == "random" {
	//			log.Printf("Study %s is completed", studyId)
	//			break
	//		}
	//		if r.checkGoal(instance, studyConfig, getMetricsReply) {
	//			log.Printf("Study %s is completed. Reach the Goal.", studyId)
	//			break
	//		}
	//	}
	//	instance.Status.Condition = katibv1alpha1.ConditionCompleted
	//	if err := r.Update(context.TODO(), instance); err != nil {
	//		log.Printf("Study: %s Condition Update error %v", studyId, err)
	//		return
	//	}
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

//func (r *ReconcileStudyJobController) getWorkerConf(wSpec *katibv1alpha1.WorkerSpec) (*katibapi.WorkerConfig, error) {
//	w := &katibapi.WorkerConfig{
//		Command: []string{},
//		Mount:   &katibapi.MountConf{},
//	}
//	if wSpec != nil {
//		w.Image = wSpec.Image
//		if wSpec.Command != nil {
//			for _, c := range wSpec.Command {
//				w.Command = append(w.Command, c)
//			}
//		}
//		w.Gpu = int32(wSpec.GPU)
//		w.Scheduler = wSpec.Scheduler
//		w.Mount.Pvc = wSpec.MountConf.Pvc
//		w.Mount.Path = wSpec.MountConf.Path
//		w.PullSecret = wSpec.PullSecret
//	}
//	return w, nil
//}

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

func (r *ReconcileStudyJobController) generate_randid() string {
	id_ := make([]byte, 2)
	_, err := rand.Read(id_)
	if err != nil {
		log.Printf("Error reading random: %v", err)
		return ""
	}
	return fmt.Sprintf("%016x", id_)[12:]
}

type WorkerInstance struct {
	WorkerId         string
	Image            string
	Command          []string
	VolumeConfigs    []katibv1alpha1.VolumeConfig
	WorkerParameters []katibv1alpha1.WorkerParameter
	HyperParameters  []katibv1alpha1.WorkerParameter
}

func (r *ReconcileStudyJobController) getWorerManifest(workerSpec *katibv1alpha1.WorkerSpec, tl []*katibapi.Trial) ([]string, []*bytes.Buffer, error) {
	wids := make([]string, len(tl))
	wins := make([]*bytes.Buffer, len(tl))

	wtp, err := template.New("DefaultWorkerTemplate").Parse(DefaultWorkerTemplate)
	if err != nil {
		return nil, nil, err
	}
	for i, t := range tl {
		wids[i] = t.TrialId + "-" + r.generate_randid()
		wi := WorkerInstance{
			WorkerId:         wids[i],
			Image:            workerSpec.Image,
			Command:          workerSpec.Command,
			VolumeConfigs:    workerSpec.VolumeConfigs,
			WorkerParameters: workerSpec.WorkerParameters,
		}
		var b bytes.Buffer
		for _, p := range t.ParameterSet {
			wi.HyperParameters = append(wi.HyperParameters,
				katibv1alpha1.WorkerParameter{
					Key:   p.Name,
					Value: p.Value,
				})
		}
		err = wtp.Execute(&b, wi)
		if err != nil {
			return nil, nil, err
		}
		wins[i] = &b
	}
	return wids, wins, nil
}

//func (r *ReconcileStudyJobController) runTrial(c katibapi.ManagerClient, studyId string, tl []*katibapi.Trial, studyConfig *katibapi.StudyConfig, wSpec *katibv1alpha1.WorkerSpec) ([]string, []katibv1alpha1.TrialSet, error) {
//	ctx := context.Background()
//	workerParameter := make(map[string][]*katibapi.Parameter)
//	workerConfig, err := r.getWorkerConf(wSpec)
//	if err != nil {
//		log.Printf("Study: %s getWorkerConf Failed %v", studyId, err)
//		return nil, nil, err
//	}
//	wl := make([]string, len(tl))
//	ts := make([]katibv1alpha1.TrialSet, len(tl))
//	for i, t := range tl {
//		wc := workerConfig
//		rtr := &katibapi.RunTrialRequest{
//			StudyId:      studyId,
//			TrialId:      t.TrialId,
//			Runtime:      "kubernetes",
//			WorkerConfig: wc,
//		}
//		for _, p := range t.ParameterSet {
//			rtr.WorkerConfig.Command = append(rtr.WorkerConfig.Command, p.Name)
//			rtr.WorkerConfig.Command = append(rtr.WorkerConfig.Command, p.Value)
//		}
//		workerReply, err := c.RunTrial(ctx, rtr)
//		if err != nil {
//			log.Printf("Study: %s RunTrial Error %v", studyId, err)
//			return nil, nil, err
//		}
//		workerParameter[workerReply.WorkerId] = t.ParameterSet
//		saveModelRequest := &katibapi.SaveModelRequest{
//			Model: &katibapi.ModelInfo{
//				StudyName:  studyConfig.Name,
//				WorkerId:   workerReply.WorkerId,
//				Parameters: t.ParameterSet,
//				Metrics:    []*katibapi.Metrics{},
//				ModelPath:  "pvc:/Path/to/Model",
//			},
//			DataSet: &katibapi.DataSetInfo{
//				Name: "Mnist",
//				Path: "/path/to/data",
//			},
//		}
//		_, err = c.SaveModel(ctx, saveModelRequest)
//		if err != nil {
//			log.Printf("Study: %s SaveModel Error %v", studyId, err)
//			return nil, nil, err
//		}
//		log.Printf("Study: %s WorkerID %s start", studyId, workerReply.WorkerId)
//		wl[i] = workerReply.WorkerId
//		ts[i].TrialId = t.TrialId
//		ts[i].WorkerIdList = append(ts[i].WorkerIdList, workerReply.WorkerId)
//	}
//	return wl, ts, nil
//}

func (r *ReconcileStudyJobController) saveOrUpdateModel(c katibapi.ManagerClient, getMetricsReply *katibapi.GetMetricsReply, studyConfig *katibapi.StudyConfig, studyId string) error {
	ctx := context.Background()
	for _, mls := range getMetricsReply.MetricsLogSets {
		if len(mls.MetricsLogs) > 0 {
			log.Printf("Study: %s WorkerID %s :", studyId, mls.WorkerId)
			//Only Metrics can be updated.
			saveModelRequest := &katibapi.SaveModelRequest{
				Model: &katibapi.ModelInfo{
					StudyName: studyConfig.Name,
					WorkerId:  mls.WorkerId,
					Metrics:   []*katibapi.Metrics{},
				},
			}
			for _, ml := range mls.MetricsLogs {
				if len(ml.Values) > 0 {
					log.Printf("\t Metrics Name %s Value %v", ml.Name, ml.Values[len(ml.Values)-1])
					saveModelRequest.Model.Metrics = append(saveModelRequest.Model.Metrics, &katibapi.Metrics{Name: ml.Name, Value: ml.Values[len(ml.Values)-1].Value})
				}
			}
			_, err := c.SaveModel(ctx, saveModelRequest)
			if err != nil {
				log.Printf("Study: %s SaveModel Error %v", studyId, err)
				return err
			}
		}
	}
	return nil
}

func (r *ReconcileStudyJobController) isCompletedAllWorker(c katibapi.ManagerClient, ms []*katibapi.MetricsLogSet) bool {
	for _, mls := range ms {
		if mls.WorkerStatus != katibapi.State_COMPLETED {
			return false
		}
	}
	return true
}
