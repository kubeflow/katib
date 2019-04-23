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

package util

import (
	//v1 "k8s.io/api/core/v1"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	trialapi "github.com/kubeflow/katib/pkg/api/v1alpha2"
)

func CreateExperimentInDB(instance *experimentsv1alpha2.Experiment) error {

	return nil
}

func UpdateExperimentStatusInDB(instance *experimentsv1alpha2.Experiment) error {

	return nil
}

func GetSuggestions(instance *experimentsv1alpha2.Experiment, addCount int) ([]*trialapi.Trial, error) {

	return nil, nil
}
