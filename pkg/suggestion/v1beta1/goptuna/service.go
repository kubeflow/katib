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

package suggestion_goptuna_v1beta1

import (
	"context"
	"sync"

	"github.com/c-bata/goptuna"
	api_v1_beta1 "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

const (
	AlgorithmCMAES  = "cmaes"
	AlgorithmTPE    = "tpe"
	AlgorithmRandom = "random"
	AlgorithmSobol  = "sobol"

	defaultStudyName = "Katib"
)

func NewSuggestionService() *SuggestionService {
	return &SuggestionService{
		searchSpace:  nil,
		study:        nil,
		trialMapping: make(map[string]int),
	}
}

type SuggestionService struct {
	mu           sync.RWMutex
	searchSpace  map[string]interface{}
	study        *goptuna.Study
	trialMapping map[string]int // Katib trial name -> Goptuna trial id
}

func (s *SuggestionService) GetSuggestions(
	ctx context.Context,
	req *api_v1_beta1.GetSuggestionsRequest,
) (*api_v1_beta1.GetSuggestionsReply, error) {
	err := s.initStudyAndSearchSpaceAtFirstRun(req.GetExperiment())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create goptuna study and search space: %s", err.Error())
	}

	objectMetricName := req.GetExperiment().GetSpec().GetObjective().GetObjectiveMetricName()
	trials, err := toGoptunaTrials(req.GetTrials(), objectMetricName, s.study, s.searchSpace)
	if err != nil {
		klog.Errorf("Failed to convert to Goptuna trials: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.syncTrials(trials)
	if err != nil {
		klog.Errorf("Failed to sync Goptuna trials: %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	currentRequestNumber := int(req.GetCurrentRequestNumber())
	parameterAssignments := make([]*api_v1_beta1.GetSuggestionsReply_ParameterAssignments, currentRequestNumber)
	for i := 0; i < currentRequestNumber; i++ {
		trialID, assignments, err := sampleNextParam(s.study, s.searchSpace)
		if err != nil {
			klog.Errorf("Failed to sample next param: trialID=%d, err=%s", trialID, err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		klog.Infof("Success to sample new trial: trialID=%d, assignments=%v", trialID, assignments)
		parameterAssignments[i] = &api_v1_beta1.GetSuggestionsReply_ParameterAssignments{
			Assignments: assignments,
		}
	}

	return &api_v1_beta1.GetSuggestionsReply{
		ParameterAssignments: parameterAssignments,
	}, nil
}

// Sync Goptuna trials with Katib trials.
func (s *SuggestionService) syncTrials(ktrials map[string]goptuna.FrozenTrial) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for katibTrialName := range ktrials {
		ktrial := ktrials[katibTrialName]
		gtrialID, found := s.trialMapping[katibTrialName]
		if !found {
			// In the CMA-ES algorithm, the parameters of Multivariate Normal Distribution MUST be updated by the
			// solutions that are sampled from the same generation. To ensure this, Goptuna stores the trial
			// metadata which contains the generation number.
			//
			// But suggestion service cannot know which Katib trial name corresponds to Goptuna trial ID.
			// Because Katib's trial name is determined by Katib controller after finished this gRPC call.
			// So `findGoptunaTrialIDByParam()` returns the goptuna trial ID from the parameter values.
			gtrialID, err = findGoptunaTrialIDByParam(s.study, s.trialMapping, ktrial)
			if err != nil {
				klog.Errorf("Failed to find Goptuna Trial ID: trialName=%s, err=%s", katibTrialName, err)
				return err
			}
			s.trialMapping[katibTrialName] = gtrialID
			klog.Infof("Update trial mapping : trialName=%s -> trialID=%d", katibTrialName, gtrialID)
		}

		gtrial, err := s.study.Storage.GetTrial(gtrialID)
		if err != nil {
			return err
		}

		// It doesn't need to update finished trials.
		if gtrial.State.IsFinished() {
			continue
		}

		if ktrial.State == gtrial.State {
			continue
		}

		if ktrial.State == goptuna.TrialStateComplete {
			err = s.study.Storage.SetTrialValue(gtrialID, ktrial.Value)
			if err != nil {
				return err
			}
		}

		err = s.study.Storage.SetTrialState(gtrialID, ktrial.State)
		if err != nil {
			klog.Errorf("Failed to update state: %s", err)
			return err
		}

		klog.Infof("Detect changes of Trial (trialName=%s, trialID=%d) : State %s, Evaluation %f",
			katibTrialName, gtrialID, ktrial.State, ktrial.Value)
	}
	return nil
}

func (s *SuggestionService) initStudyAndSearchSpaceAtFirstRun(
	experiment *api_v1_beta1.Experiment,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.study != nil && s.searchSpace != nil {
		return nil
	}

	study, searchSpace, err := createStudyAndSearchSpace(experiment)
	if err != nil {
		return err
	}

	s.study = study
	s.searchSpace = searchSpace
	return nil
}

func (s *SuggestionService) ValidateAlgorithmSettings(
	ctx context.Context,
	req *api_v1_beta1.ValidateAlgorithmSettingsRequest,
) (*api_v1_beta1.ValidateAlgorithmSettingsReply, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is empty")
	}

	algorithmName := req.GetExperiment().GetSpec().GetAlgorithm().GetAlgorithmName()
	if algorithmName != AlgorithmRandom && algorithmName != AlgorithmCMAES && algorithmName != AlgorithmTPE && algorithmName != AlgorithmSobol {
		return nil, status.Error(codes.InvalidArgument, "unsupported algorithm")
	}

	params := req.GetExperiment().GetSpec().GetParameterSpecs().GetParameters()
	if algorithmName == AlgorithmCMAES {
		cnt := 0
		for _, p := range params {
			if p.ParameterType == api_v1_beta1.ParameterType_DOUBLE || p.ParameterType == api_v1_beta1.ParameterType_INT {
				cnt++
			}
		}
		if cnt < 2 {
			return nil, status.Error(codes.InvalidArgument, "CMA-ES only supports two or more dimensional continuous search space.")
		}
	}

	paramSet := make(map[string]interface{}, len(params))
	for _, p := range params {
		if _, ok := paramSet[p.Name]; ok {
			return nil, status.Errorf(codes.InvalidArgument, "Detect duplicated parameter name: %s", p.Name)
		}
		paramSet[p.Name] = nil
	}
	_, _, err := createStudyAndSearchSpace(req.GetExperiment())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create goptuna study and search space: %s", err.Error())
	}
	return &api_v1_beta1.ValidateAlgorithmSettingsReply{}, nil
}

// This is a compile-time assertion to ensure that SuggestionService
// implements an api_v1_beta1.SuggestionServer interface.
var _ api_v1_beta1.SuggestionServer = &SuggestionService{}
