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
	"fmt"

	katibapi "github.com/kubeflow/katib/pkg/api"
)

func saveModel(c katibapi.ManagerClient, studyID string, trialID string, workerID string) error {
	ctx := context.Background()
	// Disable ModelDB
	//getStudyreq := &katibapi.GetStudyRequest{
	//	StudyId: studyId,
	//}
	//getStudyReply, err := c.GetStudy(ctx, getStudyreq)
	//if err != nil {
	//	return err
	//}
	//sc := getStudyReply.StudyConfig
	getMetricsRequest := &katibapi.GetMetricsRequest{
		StudyId:   studyID,
		WorkerIds: []string{workerID},
	}
	getMetricsReply, err := c.GetMetrics(ctx, getMetricsRequest)
	if err != nil {
		return err
	}
	for _, mls := range getMetricsReply.MetricsLogSets {
		mets := []*katibapi.Metrics{}
		var trial *katibapi.Trial = nil
		gtret, err := c.GetTrials(
			ctx,
			&katibapi.GetTrialsRequest{
				StudyId: studyID,
			})
		if err != nil {
			return err
		}
		for _, t := range gtret.Trials {
			if t.TrialId == trialID {
				trial = t
			}
		}
		for _, ml := range mls.MetricsLogs {
			if ml != nil {
				if len(ml.Values) > 0 {
					mets = append(mets, &katibapi.Metrics{
						Name:  ml.Name,
						Value: ml.Values[len(ml.Values)-1].Value,
					})
				}
			}
		}
		if trial == nil {
			return fmt.Errorf("Trial %s not found", trialID)
		}
		// Disable ModelDB
		//		if len(mets) > 0 {
		//			smr := &katibapi.SaveModelRequest{
		//				Model: &katibapi.ModelInfo{
		//					StudyName:  sc.Name,
		//					WorkerId:   mls.WorkerId,
		//					Parameters: trial.ParameterSet,
		//					Metrics:    mets,
		//					ModelPath:  sc.Name,
		//				},
		//				DataSet: &katibapi.DataSetInfo{
		//					Name: sc.Name,
		//					Path: sc.Name,
		//				},
		//			}
		//			_, err = c.SaveModel(ctx, smr)
		//			if err != nil {
		//				return err
		//			}
		//		}
	}
	return nil
}
