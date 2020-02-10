package util

import (
	"fmt"

	suggestionsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	pytorchv1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1"
	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	log = logf.Log.WithName("util-annotations")
)

// SuggestionAnnotations returns the expected suggestion annotations.
func SuggestionAnnotations(instance *suggestionsv1alpha3.Suggestion) map[string]string {
	return appendAnnotation(
		instance.Annotations,
		consts.AnnotationIstioSidecarInjectName,
		consts.AnnotationIstioSidecarInjectValue)
}

// TrainingJobAnnotations returns unstructured job with annotations.
func TrainingJobAnnotations(desiredJob *unstructured.Unstructured) error {
	kind := desiredJob.GetKind()
	switch kind {
	case consts.JobKindJob:
		batchJob := &batchv1.Job{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(desiredJob.Object, &batchJob)
		if err != nil {
			log.Error(err, "Convert unstructured to job error")
			return err
		}
		batchJob.Spec.Template.ObjectMeta.Annotations = appendAnnotation(
			batchJob.Spec.Template.ObjectMeta.Annotations,
			consts.AnnotationIstioSidecarInjectName,
			consts.AnnotationIstioSidecarInjectValue)
		desiredJob.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(batchJob)
		if err != nil {
			log.Error(err, "Convert job to unstructured error")
			return err
		}
		return nil
	case consts.JobKindTF:
		tfJob := &tfv1.TFJob{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(desiredJob.Object, &tfJob)
		if err != nil {
			log.Error(err, "Convert unstructured to TFJob error")
			return err
		}
		for _, replicaSpec := range tfJob.Spec.TFReplicaSpecs {
			replicaSpec.Template.ObjectMeta.Annotations = appendAnnotation(
				replicaSpec.Template.ObjectMeta.Annotations,
				consts.AnnotationIstioSidecarInjectName,
				consts.AnnotationIstioSidecarInjectValue)
		}
		desiredJob.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(tfJob)
		if err != nil {
			log.Error(err, "Convert TFJob to unstructured error")
			return err
		}
		return nil
	case consts.JobKindPyTorch:
		pytorchJob := &pytorchv1.PyTorchJob{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(desiredJob.Object, &pytorchJob)
		if err != nil {
			log.Error(err, "Convert unstructured to PytorchJob error")
			return err
		}
		for _, replicaSpec := range pytorchJob.Spec.PyTorchReplicaSpecs {
			replicaSpec.Template.ObjectMeta.Annotations = appendAnnotation(
				replicaSpec.Template.ObjectMeta.Annotations,
				consts.AnnotationIstioSidecarInjectName,
				consts.AnnotationIstioSidecarInjectValue)
		}
		desiredJob.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(pytorchJob)
		if err != nil {
			log.Error(err, "Convert PytorchJob to unstructured error")
			return err
		}
		return nil
	default:
		return fmt.Errorf("Invalid Katib Training Job kind %v", kind)
	}

}

func appendAnnotation(annotations map[string]string, newAnnotationName string, newAnnotationValue string) map[string]string {
	res := make(map[string]string)
	for k, v := range annotations {
		res[k] = v
	}
	res[newAnnotationName] = newAnnotationValue

	return res
}
