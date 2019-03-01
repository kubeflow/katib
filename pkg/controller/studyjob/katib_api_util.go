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

package studyjob

import (
	"context"
	"log"

	"github.com/kubeflow/katib/pkg"
	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"google.golang.org/grpc"
)

func initializeStudy(instance *katibv1alpha1.StudyJob) error {
	if instance.Spec.SuggestionSpec.SuggestionAlgorithm == "" {
		instance.Spec.SuggestionSpec.SuggestionAlgorithm = "random"
	}

	instance.Status.Condition = katibv1alpha1.ConditionRunning

	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect katib manager error %v", err)
		instance.Status.Condition = katibv1alpha1.ConditionFailed
		return nil
	}
	defer conn.Close()
	c := katibapi.NewManagerClient(conn)

	studyConfig, err := getStudyConf(instance)
	if err != nil {
		return err
	}
	log.Printf("Create Study %s", studyConfig.Name)
	//CreateStudy
	studyID, err := createStudy(c, studyConfig)
	if err != nil {
		return err
	}
	instance.Status.StudyID = studyID
	log.Printf("Study: %s Suggestion Spec %v", studyID, instance.Spec.SuggestionSpec)
	var sspec *katibv1alpha1.SuggestionSpec
	if instance.Spec.SuggestionSpec != nil {
		sspec = instance.Spec.SuggestionSpec
	} else {
		sspec = &katibv1alpha1.SuggestionSpec{}
	}
	sspec.SuggestionParameters = append(sspec.SuggestionParameters,
		katibapi.SuggestionParameter{
			Name:  "SuggestionCount",
			Value: "0",
		})
	sPID, err := setSuggestionParam(c, studyID, sspec)
	if err != nil {
		return err
	}
	instance.Status.SuggestionParameterID = sPID
	instance.Status.Condition = katibv1alpha1.ConditionRunning
	return nil
}

func getStudyConf(instance *katibv1alpha1.StudyJob) (*katibapi.StudyConfig, error) {
	jobType := getJobType(instance)
	if jobType == jobTypeNAS {
		return populateConfigForNAS(instance)
	}
	return populateConfigForHP(instance)
}

func getJobType(instance *katibv1alpha1.StudyJob) string {
	if instance.Spec.NasConfig != nil {
		return jobTypeNAS
	}
	return jobTypeHP
}

func populateCommonConfigFields(instance *katibv1alpha1.StudyJob, sconf *katibapi.StudyConfig) {

	sconf.Name = instance.Spec.StudyName
	sconf.Owner = instance.Spec.Owner
	if instance.Spec.OptimizationGoal != nil {
		sconf.OptimizationGoal = *instance.Spec.OptimizationGoal
	}
	sconf.ObjectiveValueName = instance.Spec.ObjectiveValueName
	switch instance.Spec.OptimizationType {
	case katibv1alpha1.OptimizationTypeMinimize:
		sconf.OptimizationType = katibapi.OptimizationType_MINIMIZE
	case katibv1alpha1.OptimizationTypeMaximize:
		sconf.OptimizationType = katibapi.OptimizationType_MAXIMIZE
	default:
		sconf.OptimizationType = katibapi.OptimizationType_UNKNOWN_OPTIMIZATION
	}
	for _, m := range instance.Spec.MetricsNames {
		sconf.Metrics = append(sconf.Metrics, m)
	}
	sconf.JobId = string(instance.UID)
}

func populateConfigForHP(instance *katibv1alpha1.StudyJob) (*katibapi.StudyConfig, error) {
	sconf := &katibapi.StudyConfig{
		Metrics: []string{},
		ParameterConfigs: &katibapi.StudyConfig_ParameterConfigs{
			Configs: []*katibapi.ParameterConfig{},
		},
	}

	populateCommonConfigFields(instance, sconf)

	for _, pc := range instance.Spec.ParameterConfigs {
		p := &katibapi.ParameterConfig{
			Feasible: &katibapi.FeasibleSpace{},
		}
		p.Name = pc.Name
		if pc.Feasible.Min != "" && pc.Feasible.Max != "" {
			p.Feasible.Min = pc.Feasible.Min
			p.Feasible.Max = pc.Feasible.Max
		}
		if pc.Feasible.List != nil {
			p.Feasible.List = pc.Feasible.List
		}

		if pc.Feasible.Step != "" {
			p.Feasible.Step = pc.Feasible.Step
		}
		switch pc.ParameterType {
		case katibv1alpha1.ParameterTypeUnknown:
			p.ParameterType = katibapi.ParameterType_UNKNOWN_TYPE
		case katibv1alpha1.ParameterTypeDouble:
			p.ParameterType = katibapi.ParameterType_DOUBLE
		case katibv1alpha1.ParameterTypeInt:
			p.ParameterType = katibapi.ParameterType_INT
		case katibv1alpha1.ParameterTypeDiscrete:
			p.ParameterType = katibapi.ParameterType_DISCRETE
		case katibv1alpha1.ParameterTypeCategorical:
			p.ParameterType = katibapi.ParameterType_CATEGORICAL
		}
		sconf.ParameterConfigs.Configs = append(sconf.ParameterConfigs.Configs, p)
	}

	sconf.JobType = jobTypeHP
	return sconf, nil
}

func populateConfigForNAS(instance *katibv1alpha1.StudyJob) (*katibapi.StudyConfig, error) {
	sconf := &katibapi.StudyConfig{
		Metrics: []string{},
		NasConfig: &katibapi.NasConfig{
			GraphConfig: &katibapi.GraphConfig{},
			Operations: &katibapi.NasConfig_Operations{
				Operation: []*katibapi.Operation{},
			},
		},
	}
	populateCommonConfigFields(instance, sconf)

	sconf.NasConfig.GraphConfig.NumLayers = instance.Spec.NasConfig.GraphConfig.NumLayers
	for _, i := range instance.Spec.NasConfig.GraphConfig.InputSize {
		sconf.NasConfig.GraphConfig.InputSize = append(sconf.NasConfig.GraphConfig.InputSize, i)
	}
	for _, o := range instance.Spec.NasConfig.GraphConfig.OutputSize {
		sconf.NasConfig.GraphConfig.OutputSize = append(sconf.NasConfig.GraphConfig.OutputSize, o)
	}
	for _, op := range instance.Spec.NasConfig.Operations {
		operation := &katibapi.Operation{
			ParameterConfigs: &katibapi.Operation_ParameterConfigs{
				Configs: []*katibapi.ParameterConfig{},
			},
		}
		operation.OperationType = op.OperationType
		for _, pc := range op.ParameterConfigs {
			p := &katibapi.ParameterConfig{
				Feasible: &katibapi.FeasibleSpace{},
			}

			p.Name = pc.Name
			if pc.Feasible.Min != "" && pc.Feasible.Max != "" {
				p.Feasible.Min = pc.Feasible.Min
				p.Feasible.Max = pc.Feasible.Max
			}
			if pc.Feasible.List != nil {

				p.Feasible.List = pc.Feasible.List

			}
			if pc.Feasible.Step != "" {
				p.Feasible.Step = pc.Feasible.Step
			}
			switch pc.ParameterType {
			case katibv1alpha1.ParameterTypeUnknown:
				p.ParameterType = katibapi.ParameterType_UNKNOWN_TYPE
			case katibv1alpha1.ParameterTypeDouble:
				p.ParameterType = katibapi.ParameterType_DOUBLE
			case katibv1alpha1.ParameterTypeInt:
				p.ParameterType = katibapi.ParameterType_INT
			case katibv1alpha1.ParameterTypeDiscrete:
				p.ParameterType = katibapi.ParameterType_DISCRETE
			case katibv1alpha1.ParameterTypeCategorical:
				p.ParameterType = katibapi.ParameterType_CATEGORICAL
			}

			operation.ParameterConfigs.Configs = append(operation.ParameterConfigs.Configs, p)
		}
		sconf.NasConfig.Operations.Operation = append(sconf.NasConfig.Operations.Operation, operation)
	}

	sconf.JobType = jobTypeNAS
	return sconf, nil
}

func deleteStudy(instance *katibv1alpha1.StudyJob) error {
	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Connect katib manager error %v", err)
		return err
	}
	defer conn.Close()
	c := katibapi.NewManagerClient(conn)
	ctx := context.Background()
	studyID := instance.Status.StudyID
	if studyID == "" {
		// in case that information for a studyjob is not created in DB
		return nil
	}
	deleteStudyreq := &katibapi.DeleteStudyRequest{
		StudyId: studyID,
	}
	if _, err = c.DeleteStudy(ctx, deleteStudyreq); err != nil {
		log.Printf("DeleteStudy error %v", err)
		return err
	}
	return nil
}

func createStudy(c katibapi.ManagerClient, studyConfig *katibapi.StudyConfig) (string, error) {
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
	log.Printf("Study ID %s StudyConf %v", studyID, getStudyReply.StudyConfig)
	return studyID, nil
}

func setSuggestionParam(c katibapi.ManagerClient, studyID string, suggestionSpec *katibv1alpha1.SuggestionSpec) (string, error) {
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

func getSuggestionParam(c katibapi.ManagerClient, paramID string) ([]*katibapi.SuggestionParameter, error) {
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

func getSuggestion(c katibapi.ManagerClient, studyID string, suggestionSpec *katibv1alpha1.SuggestionSpec, sParamID string) (*katibapi.GetSuggestionsReply, error) {
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
