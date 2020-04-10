package suggestion_goptuna_v1alpha3

import (
	"context"

	"github.com/c-bata/goptuna"
	"github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

const (
	AlgorithmCMAES  = "cmaes"
	AlgorithmTPE    = "tpe"
	AlgorithmRandom = "random"

	defaultStudyName = "Katib"
)

func NewSuggestionService() *SuggestionService {
	return &SuggestionService{}
}

type SuggestionService struct{}

func (s *SuggestionService) GetSuggestions(
	ctx context.Context,
	req *api_v1_alpha3.GetSuggestionsRequest,
) (*api_v1_alpha3.GetSuggestionsReply, error) {
	direction, err := toGoptunaDirection(req.GetExperiment().GetSpec().GetObjective().GetType())
	if err != nil {
		klog.Errorf("Failed to convert to Goptuna direction: %s", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	independentSampler, relativeSampler, err := toGoptunaSampler(req.GetExperiment().GetSpec().GetAlgorithm())
	if err != nil {
		klog.Errorf("Failed to create Goptuna sampler: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	searchSpace, err := toGoptunaSearchSpace(req.GetExperiment().GetSpec().GetParameterSpecs().GetParameters())
	if err != nil {
		klog.Errorf("Failed to convert to Goptuna search space: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("Goptuna search space: %#v", searchSpace)

	studyOpts := make([]goptuna.StudyOption, 0, 3)
	studyOpts = append(studyOpts, goptuna.StudyOptionSetDirection(direction))
	if independentSampler != nil {
		studyOpts = append(studyOpts, goptuna.StudyOptionSampler(independentSampler))
	}
	if relativeSampler != nil {
		studyOpts = append(studyOpts, goptuna.StudyOptionRelativeSampler(relativeSampler))
	}

	study, err := goptuna.CreateStudy(defaultStudyName, studyOpts...)
	if err != nil {
		klog.Errorf("Failed to create Goptuna study: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	objectMetricName := req.GetExperiment().GetSpec().GetObjective().GetObjectiveMetricName()
	trials, err := toGoptunaTrials(req.GetTrials(), objectMetricName, study, searchSpace)
	if err != nil {
		klog.Errorf("Failed to convert to Goptuna trials: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	for _, t := range trials {
		_, err = study.Storage.CloneTrial(study.ID, t)
		if err != nil {
			klog.Errorf("Failed to register trials: %s", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	requestNumber := int(req.GetRequestNumber())
	parameterAssignments := make([]*api_v1_alpha3.GetSuggestionsReply_ParameterAssignments, requestNumber)
	for i := 0; i < requestNumber; i++ {
		assignments, err := sampleNextParam(study, searchSpace)
		if err != nil {
			klog.Errorf("Failed to sample next param: %s", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		parameterAssignments[i] = &api_v1_alpha3.GetSuggestionsReply_ParameterAssignments{
			Assignments: assignments,
		}
	}

	klog.Infof("Success to sample %d parameters", requestNumber)
	return &api_v1_alpha3.GetSuggestionsReply{
		ParameterAssignments: parameterAssignments,
		Algorithm: &api_v1_alpha3.AlgorithmSpec{
			AlgorithmName:     "",
			AlgorithmSetting:  nil,
			EarlyStoppingSpec: &api_v1_alpha3.EarlyStoppingSpec{},
		},
	}, nil
}

func (s *SuggestionService) ValidateAlgorithmSettings(ctx context.Context, req *api_v1_alpha3.ValidateAlgorithmSettingsRequest) (*api_v1_alpha3.ValidateAlgorithmSettingsReply, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is empty")
	}

	algorithmName := req.GetExperiment().GetSpec().GetAlgorithm().GetAlgorithmName()
	if algorithmName != AlgorithmRandom && algorithmName != AlgorithmCMAES && algorithmName != AlgorithmTPE {
		return nil, status.Error(codes.InvalidArgument, "unsupported algorithm")
	}

	params := req.GetExperiment().GetSpec().GetParameterSpecs().GetParameters()
	if algorithmName == AlgorithmCMAES && len(params) < 2 {
		return nil, status.Error(codes.InvalidArgument, "CMA-ES only supports two or more dimensional continuous search space.")
	}

	paramSet := make(map[string]interface{}, len(params))
	for _, p := range params {
		if _, ok := paramSet[p.Name]; ok {
			return nil, status.Errorf(codes.InvalidArgument, "Detect duplicated parameter name: %s", p.Name)
		}
		paramSet[p.Name] = nil
	}

	return &api_v1_alpha3.ValidateAlgorithmSettingsReply{}, nil
}

// This is a compile-time assertion to ensure that SuggestionService
// implements an api_v1_alpha3.SuggestionServer interface.
var _ api_v1_alpha3.SuggestionServer = &SuggestionService{}
