package validator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	batchv1beta "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	commonapiv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	commonv1alpha3 "github.com/kubeflow/katib/pkg/common/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/managerclient"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/experiment/manifest"
)

var log = logf.Log.WithName("experiment-validating-webhook")

type Validator interface {
	ValidateExperiment(instance *experimentsv1alpha3.Experiment) error
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

func (g *DefaultValidator) ValidateExperiment(instance *experimentsv1alpha3.Experiment) error {
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

func (g *DefaultValidator) validateAlgorithmSettings(inst *experimentsv1alpha3.Experiment) error {

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

func (g *DefaultValidator) validateObjective(obj *commonapiv1alpha3.ObjectiveSpec) error {
	if obj == nil {
		return fmt.Errorf("No spec.objective specified.")
	}
	if obj.Type != commonapiv1alpha3.ObjectiveTypeMinimize && obj.Type != commonapiv1alpha3.ObjectiveTypeMaximize {
		return fmt.Errorf("spec.objective.type must be %s or %s.", commonapiv1alpha3.ObjectiveTypeMinimize, commonapiv1alpha3.ObjectiveTypeMaximize)
	}
	if obj.ObjectiveMetricName == "" {
		return fmt.Errorf("No spec.objective.objectiveMetricName specified.")
	}
	return nil
}

func (g *DefaultValidator) validateAlgorithm(ag *commonapiv1alpha3.AlgorithmSpec) error {
	if ag == nil {
		return fmt.Errorf("No spec.algorithm specified.")
	}
	if ag.AlgorithmName == "" {
		return fmt.Errorf("No spec.algorithm.name specified.")
	}

	return nil
}

func (g *DefaultValidator) validateTrialTemplate(instance *experimentsv1alpha3.Experiment) error {
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
	supportedJobs := commonv1alpha3.GetSupportedJobList()
	for _, sJob := range supportedJobs {
		if gvk == sJob {
			return nil
		}
	}
	return fmt.Errorf("Job type %v not supported", gvk)
}

func (g *DefaultValidator) validateForCreate(inst *experimentsv1alpha3.Experiment) error {
	reply, err := g.PreCheckRegisterExperimentInDB(inst)
	if err != nil {
		return fmt.Errorf("Fail to check record for the experiment in DB: %v", err)
	} else if !reply.CanRegister {
		return fmt.Errorf("Record for the experiment has existed in DB; Please try to rename the experiment")
	} else {
		return nil
	}
}

func (g *DefaultValidator) validateMetricsCollector(inst *experimentsv1alpha3.Experiment) error {
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

	// Above part of this method will be dropped
	mcSpec := inst.Spec.MetricsCollectorSpec
	mcKind := mcSpec.Collector.Kind
	// TODO(hougangliu):
	// 1. validate .spec.metricsCollectorSpec.source.filter
	// 2. log warning message if some field will not be used for the metricsCollector kind
	switch mcKind {
	case commonapiv1alpha3.NoneCollector, commonapiv1alpha3.StdOutCollector:
		return nil
	case commonapiv1alpha3.FileCollector:
		if mcSpec.Source == nil || mcSpec.Source.FileSystemPath == nil ||
			mcSpec.Source.FileSystemPath.Kind != commonapiv1alpha3.FileKind || !filepath.IsAbs(mcSpec.Source.FileSystemPath.Path) {
			return fmt.Errorf("File path where metrics file exists is required by .spec.metricsCollectorSpec.source.fileSystemPath.path")
		}
	case commonapiv1alpha3.TfEventCollector:
		if mcSpec.Source == nil || mcSpec.Source.FileSystemPath == nil ||
			mcSpec.Source.FileSystemPath.Kind != commonapiv1alpha3.DirectoryKind || !filepath.IsAbs(mcSpec.Source.FileSystemPath.Path) {
			return fmt.Errorf("Directory path where tensorflow event files exist is required by .spec.metricsCollectorSpec.source.fileSystemPath.path")
		}
	case commonapiv1alpha3.PrometheusMetricCollector:
		i, err := strconv.Atoi(mcSpec.Source.HttpGet.Port.String())
		if err != nil || i <= 0 {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.httpGet.port must be a positive integer value for metrics collector kind: %v.", mcKind)
		}
		if !strings.HasPrefix(mcSpec.Source.HttpGet.Path, "/") {
			return fmt.Errorf(".spec.metricsCollectorSpec.source.httpGet.path is invalid for metrics collector kind: %v.", mcKind)
		}
	case commonapiv1alpha3.CustomCollector:
		if mcSpec.Collector.CustomCollector == nil {
			return fmt.Errorf(".spec.metricsCollectorSpec.collector.customCollector is required for metrics collector kind: %v.", mcKind)
		}
	default:
		return fmt.Errorf("Invalid metrics collector kind: %v.", mcKind)
	}

	return nil
}
