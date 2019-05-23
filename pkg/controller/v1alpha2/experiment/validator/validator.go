package validator

import (
	"bytes"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	batchv1beta "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/common/v1alpha2"
	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	commonv1alpha2 "github.com/kubeflow/katib/pkg/common/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/experiment/managerclient"
	"github.com/kubeflow/katib/pkg/controller/v1alpha2/experiment/manifest"
)

var log = logf.Log.WithName("experiment-controller")

type Validator interface {
	ValidateExperiment(instance *experimentsv1alpha2.Experiment) error
	InjectClient(c client.Client)
}

type DefaultValidator struct {
	manifest.Generator
	managerclient.ManagerClient
}

func New(generator manifest.Generator, managerClient managerclient.ManagerClient) Validator {
	return &DefaultValidator{
		Generator:     generator,
		ManagerClient: managerClient,
	}
}

func (g *DefaultValidator) InjectClient(c client.Client) {
	g.Generator.InjectClient(c)
}

func (g *DefaultValidator) ValidateExperiment(instance *experimentsv1alpha2.Experiment) error {
	if !instance.IsCreated() {
		if err := g.validateForCreate(instance); err != nil {
			return err
		}
	}
	if err := g.validateObjective(instance.Spec.Objective); err != nil {
		return err
	}
	if err := g.validateAlgorithm(instance.Spec.Algorithm); err != nil {
		return err
	}

	if err := g.validateTrialTemplate(instance); err != nil {
		return err
	}

	if len(instance.Spec.Parameters) == 0 && instance.Spec.NasConfig == nil {
		return fmt.Errorf("spec.parameters or spec.nasConfig must be specified.")
	}

	if len(instance.Spec.Parameters) > 0 && instance.Spec.NasConfig != nil {
		return fmt.Errorf("Only one of spec.parameters and spec.nasConfig can be specified.")
	}

	if err := g.validateAlgorithmSettings(instance); err != nil {
		return err
	}

	if err := g.validateMetricsCollector(instance); err != nil {
		return err
	}

	return nil
}

func (g *DefaultValidator) validateAlgorithmSettings(inst *experimentsv1alpha2.Experiment) error {

	_, err := g.ValidateAlgorithmSettings(inst)
	statusCode, _ := status.FromError(err)

	if statusCode.Code() == codes.Unknown {
		return fmt.Errorf("Method ValidateAlgorithmSettings not found inside Suggestion service: %s", inst.Spec.Algorithm.AlgorithmName)
	}

	if statusCode.Code() == codes.InvalidArgument || statusCode.Code() == codes.Unavailable {
		return fmt.Errorf("ValidateAlgorithmSettings Error: %v", statusCode.Message())
	}
	return nil

}

func (g *DefaultValidator) validateObjective(obj *commonapiv1alpha2.ObjectiveSpec) error {
	if obj == nil {
		return fmt.Errorf("No spec.objective specified.")
	}
	if obj.Type != commonapiv1alpha2.ObjectiveTypeMinimize && obj.Type != commonapiv1alpha2.ObjectiveTypeMaximize {
		return fmt.Errorf("spec.objective.type must be %s or %s.", commonapiv1alpha2.ObjectiveTypeMinimize, commonapiv1alpha2.ObjectiveTypeMaximize)
	}
	if obj.ObjectiveMetricName == "" {
		return fmt.Errorf("No spec.objective.objectiveMetricName specified.")
	}
	return nil
}

func (g *DefaultValidator) validateAlgorithm(ag *experimentsv1alpha2.AlgorithmSpec) error {
	if ag == nil {
		return fmt.Errorf("No spec.algorithm specified.")
	}
	if ag.AlgorithmName == "" {
		return fmt.Errorf("No spec.algorithm.name specified.")
	}

	return nil
}

func (g *DefaultValidator) validateTrialTemplate(instance *experimentsv1alpha2.Experiment) error {
	trialName := fmt.Sprintf("%s-trial", instance.GetName())
	runSpec, err := g.GetRunSpec(instance, instance.GetName(), trialName, instance.GetNamespace())
	if err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	bufSize := 1024
	buf := bytes.NewBufferString(runSpec)

	job := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, bufSize).Decode(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	if err := g.validateSupportedJob(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	if job.GetNamespace() != instance.GetNamespace() {
		return fmt.Errorf("Invalid spec.trialTemplate: metadata.namespace should be %s or {{.NameSpace}}", instance.GetNamespace())
	}
	if job.GetName() != trialName {
		return fmt.Errorf("Invalid spec.trialTemplate: metadata.name should be {{.Trial}}")
	}
	return nil
}

func (g *DefaultValidator) validateSupportedJob(job *unstructured.Unstructured) error {
	gvk := job.GroupVersionKind()
	supportedJobs := commonv1alpha2.GetSupportedJobList()
	for _, sJob := range supportedJobs {
		if gvk == sJob {
			return nil
		}
	}
	return fmt.Errorf("Job type %v not supported", gvk)
}

func (g *DefaultValidator) validateForCreate(inst *experimentsv1alpha2.Experiment) error {
	reply, err := g.PreCheckRegisterExperimentInDB(inst)
	if err != nil {
		return fmt.Errorf("Fail to check record for the experiment in DB: %v", err)
	} else if !reply.CanRegister {
		return fmt.Errorf("Record for the experiment has existed in DB; Please try to rename the experiment")
	} else {
		return nil
	}
}

func (g *DefaultValidator) validateMetricsCollector(inst *experimentsv1alpha2.Experiment) error {
	BUFSIZE := 1024
	experimentName := inst.GetName()
	trialName := fmt.Sprintf("%s-trial", inst.GetName())
	namespace := inst.GetNamespace()
	var metricNames []string
	metricNames = append(metricNames, inst.Spec.Objective.ObjectiveMetricName)
	for _, mn := range inst.Spec.Objective.AdditionalMetricNames {
		metricNames = append(metricNames, mn)
	}

	runSpec, err := g.GetRunSpec(inst, experimentName, trialName, namespace)
	if err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	buf := bytes.NewBufferString(runSpec)

	job := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, BUFSIZE).Decode(job); err != nil {
		return fmt.Errorf("Invalid spec.trialTemplate: %v.", err)
	}

	mcm, err := g.GetMetricsCollectorManifest(experimentName, trialName, job.GetKind(), namespace, metricNames, inst.Spec.MetricsCollectorSpec)
	if err != nil {
		log.Info("getMetricsCollectorManifest error", "err", err)
		return err
	}

	buf = bytes.NewBufferString(mcm)

	mcjob := &batchv1beta.CronJob{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, BUFSIZE).Decode(&mcjob); err != nil {
		log.Info("MetricsCollector Yaml decode error", "err", err)
		return err
	}

	if mcjob.GetNamespace() != namespace || mcjob.GetName() != trialName {
		return fmt.Errorf("Invalid metricsCollector template.")
	}
	return nil
}
