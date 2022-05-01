/*
Copyright 2022 The Kubeflow Authors.

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

package suggestionclient

import (
	experimentsv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
)

func convertNasConfig(nasConfig *experimentsv1beta1.NasConfig) *suggestionapi.NasConfig {
	res := &suggestionapi.NasConfig{
		GraphConfig: convertGraphConfig(nasConfig.GraphConfig),
		Operations:  convertOperations(nasConfig.Operations),
	}
	return res
}

func convertGraphConfig(graphConfig experimentsv1beta1.GraphConfig) *suggestionapi.GraphConfig {
	gc := &suggestionapi.GraphConfig{}
	if graphConfig.NumLayers != nil {
		gc.NumLayers = *graphConfig.NumLayers
	}
	gc.InputSizes = graphConfig.InputSizes
	gc.OutputSizes = graphConfig.OutputSizes
	return gc
}

func convertOperations(operations []experimentsv1beta1.Operation) *suggestionapi.NasConfig_Operations {
	ops := &suggestionapi.NasConfig_Operations{
		Operation: make([]*suggestionapi.Operation, 0),
	}
	for _, operation := range operations {
		op := &suggestionapi.Operation{
			OperationType:  operation.OperationType,
			ParameterSpecs: convertNasParameterSpecs(operation.Parameters),
		}
		ops.Operation = append(ops.Operation, op)
	}
	return ops
}

func convertNasParameterSpecs(parameters []experimentsv1beta1.ParameterSpec) *suggestionapi.Operation_ParameterSpecs {
	ps := &suggestionapi.Operation_ParameterSpecs{
		Parameters: convertParameters(parameters),
	}
	return ps
}
