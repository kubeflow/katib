// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package studyjobcontroller

import (
	"context"
	"log"

	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
)

func (r *ReconcileStudyJobController) createStudy(c katibapi.ManagerClient, studyConfig *katibapi.StudyConfig) (string, error) {
	ctx := context.Background()
	createStudyreq := &katibapi.CreateStudyRequest{
		StudyConfig: studyConfig,
	}
	createStudyreply, err := c.CreateStudy(ctx, createStudyreq)
	if err != nil {
		log.Printf("CreateStudy Error %v", err)
		return "", err
	}
	studyID := createStudyreply.StudyId
	log.Printf("Study ID %s", studyID)
	getStudyreq := &katibapi.GetStudyRequest{
		StudyId: studyID,
	}
	getStudyReply, err := c.GetStudy(ctx, getStudyreq)
	if err != nil {
		log.Printf("Study: %s GetConfig Error %v", studyID, err)
		return "", err
	}
	log.Printf("Study ID %s StudyConf%v", studyID, getStudyReply.StudyConfig)
	return studyID, nil
}

func (r *ReconcileStudyJobController) setSuggestionParam(c katibapi.ManagerClient, studyID string, suggestionSpec *katibv1alpha1.SuggestionSpec) (string, error) {
	ctx := context.Background()
	pid := ""
	if suggestionSpec.SuggestionParameters != nil {
		sspr := &katibapi.SetSuggestionParametersRequest{
			StudyId:             studyID,
			SuggestionAlgorithm: suggestionSpec.SuggestionAlgorithm,
		}
		for _, p := range suggestionSpec.SuggestionParameters {
			sspr.SuggestionParameters = append(
				sspr.SuggestionParameters,
				&katibapi.SuggestionParameter{
					Name:  p.Name,
					Value: p.Value,
				},
			)
		}
		setSuggesitonParameterReply, err := c.SetSuggestionParameters(ctx, sspr)
		if err != nil {
			log.Printf("Study %s SetConfig Error %v", studyID, err)
			return "", err
		}
		log.Printf("Study: %s setSuggesitonParameterReply %v", studyID, setSuggesitonParameterReply)
		pid = setSuggesitonParameterReply.ParamId
	}
	return pid, nil
}

func (r *ReconcileStudyJobController) getSuggestionParam(c katibapi.ManagerClient, paramID string) ([]*katibapi.SuggestionParameter, error) {
	ctx := context.Background()
	gsreq := &katibapi.GetSuggestionParametersRequest{
		ParamId: paramID,
	}
	gsrep, err := c.GetSuggestionParameters(ctx, gsreq)
	if err != nil {
		return nil, err
	}
	return gsrep.SuggestionParameters, err
}
func (r *ReconcileStudyJobController) getSuggestion(c katibapi.ManagerClient, studyID string, suggestionSpec *katibv1alpha1.SuggestionSpec, sParamID string) (*katibapi.GetSuggestionsReply, error) {
	ctx := context.Background()
	getSuggestRequest := &katibapi.GetSuggestionsRequest{
		StudyId:             studyID,
		SuggestionAlgorithm: suggestionSpec.SuggestionAlgorithm,
		RequestNumber:       int32(suggestionSpec.RequestNumber),
		//RequestNumber=0 means get all grids.
		ParamId: sParamID,
	}
	getSuggestReply, err := c.GetSuggestions(ctx, getSuggestRequest)
	if err != nil {
		log.Printf("Study: %s GetSuggestion Error %v", studyID, err)
		return nil, err
	}
	log.Printf("Study: %s CreatedTrials :", studyID)
	for _, t := range getSuggestReply.Trials {
		log.Printf("\t%v", t)
	}
	return getSuggestReply, nil
}

