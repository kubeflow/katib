package suggestionclient

import (
	experimentsv1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	suggestionapi "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
)

func convertNasConfig(nasConfig *experimentsv1alpha3.NasConfig) *suggestionapi.NasConfig {
	res := &suggestionapi.NasConfig{
		GraphConfig: convertGraphConfig(nasConfig.GraphConfig),
		Operations:  convertOperations(nasConfig.Operations),
	}
	return res
}

func convertGraphConfig(graphConfig experimentsv1alpha3.GraphConfig) *suggestionapi.GraphConfig {
	gc := &suggestionapi.GraphConfig{}
	if graphConfig.NumLayers != nil {
		gc.NumLayers = *graphConfig.NumLayers
	}
	gc.InputSizes = graphConfig.InputSizes
	gc.OutputSizes = graphConfig.OutputSizes
	return gc
}

func convertOperations(operations []experimentsv1alpha3.Operation) *suggestionapi.NasConfig_Operations {
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

func convertNasParameterSpecs(parameters []experimentsv1alpha3.ParameterSpec) *suggestionapi.Operation_ParameterSpecs {
	ps := &suggestionapi.Operation_ParameterSpecs{
		Parameters: convertParameters(parameters),
	}
	return ps
}
