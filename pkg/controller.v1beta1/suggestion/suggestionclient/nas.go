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
