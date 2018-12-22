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
package studyjob

import (
	"fmt"
	"log"

	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	pytorchjobv1beta1 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1beta1"
	tfjobv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta1"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func createWorkerJobObj(kind string) runtime.Object {
	switch kind {
	case DefaultJobWorker:
		return &batchv1.Job{}
	case TFJobWorker:
		return &tfjobv1beta1.TFJob{}
	case PyTorchJobWorker:
		return &pytorchjobv1beta1.PyTorchJob{}
	}
	return nil
}

func validateStudy(instance *katibv1alpha1.StudyJob, namespace string) error {
	if instance.Spec.SuggestionSpec == nil {
		return fmt.Errorf("No Spec.SuggestionSpec specified.")
	}
	BUFSIZE := 1024
	wkind, err := getWorkerKind(instance.Spec.WorkerSpec)
	if err != nil {
		log.Printf("getWorkerKind error %v", err)
		return err
	}

	studyID := "studyID4Validation"
	trialID := "trialID4Validation"
	workerID, wm, err := getWorkerManifest(
		nil,
		studyID,
		&katibapi.Trial{
			TrialId:      trialID,
			ParameterSet: []*katibapi.Parameter{},
		},
		instance.Spec.WorkerSpec,
		wkind,
		namespace,
		true,
	)
	if err != nil {
		return err
	}

	job := createWorkerJobObj(wkind)
	if err := k8syaml.NewYAMLOrJSONDecoder(wm, BUFSIZE).Decode(job); err != nil {
		log.Printf("Yaml decode error %v", err)
		return err
	}

	metav1Job := job.(metav1.Object)
	if metav1Job.GetNamespace() != namespace || metav1Job.GetName() != workerID {
		return fmt.Errorf("Invalid worker template.")
	}

	var mcjob batchv1beta.CronJob
	mcm, err := getMetricsCollectorManifest(studyID, trialID, workerID, wkind, namespace, instance.Spec.MetricsCollectorSpec)
	if err != nil {
		log.Printf("getMetricsCollectorManifest error %v", err)
		return err
	}

	if err := k8syaml.NewYAMLOrJSONDecoder(mcm, BUFSIZE).Decode(&mcjob); err != nil {
		log.Printf("MetricsCollector Yaml decode error %v", err)
		return err
	}

	if mcjob.GetNamespace() != namespace || mcjob.GetName() != workerID {
		return fmt.Errorf("Invalid metricsCollector template.")
	}
	return nil
}
