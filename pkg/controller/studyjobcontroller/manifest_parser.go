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
	"bytes"
	"context"
	"text/template"

	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
)

func (r *ReconcileStudyJobController) getWorkerManifest(c katibapi.ManagerClient, studyID string, trial *katibapi.Trial, workerSpec *katibv1alpha1.WorkerSpec, kind string, dryrun bool) (string, *bytes.Buffer, error) {
	var wtp *template.Template = nil
	var err error
	if workerSpec != nil {
		if workerSpec.GoTemplate.RawTemplate != "" {
			wtp, err = template.New("Worker").Parse(workerSpec.GoTemplate.RawTemplate)
		} else if workerSpec.GoTemplate.TemplatePath != "" {
			wtp, err = template.ParseFiles(workerSpec.GoTemplate.TemplatePath)
		}
		if err != nil {
			return "", nil, err
		}
	}
	if wtp == nil {
		wtp, err = template.ParseFiles("/worker-template/defaultWorkerTemplate.yaml")
		if err != nil {
			return "", nil, err
		}
	}
	var wid string
	if dryrun {
		wid = "validation"
	} else {
		cwreq := &katibapi.RegisterWorkerRequest{
			Worker: &katibapi.Worker{
				StudyId: studyID,
				TrialId: trial.TrialId,
				Status:  katibapi.State_PENDING,
				Type:    kind,
			},
		}
		cwrep, err := c.RegisterWorker(context.Background(), cwreq)
		if err != nil {
			return "", nil, err
		}
		wid = cwrep.WorkerId
	}

	wi := WorkerInstance{
		StudyID:  studyID,
		TrialID:  trial.TrialId,
		WorkerID: wid,
	}
	var b bytes.Buffer
	for _, p := range trial.ParameterSet {
		wi.HyperParameters = append(wi.HyperParameters, p)
	}
	err = wtp.Execute(&b, wi)
	if err != nil {
		return "", nil, err
	}
	return wid, &b, nil
}

func (r *ReconcileStudyJobController) getMetricsCollectorManifest(studyID string, trialID string, workerID string, namespace string, mcs *katibv1alpha1.MetricsCollectorSpec) (*bytes.Buffer, error) {
	var mtp *template.Template = nil
	var err error
	tmpValues := map[string]string{
		"StudyID":   studyID,
		"TrialID":   trialID,
		"WorkerID":  workerID,
		"NameSpace": namespace,
	}
	if mcs != nil {
		if mcs.GoTemplate.RawTemplate != "" {
			mtp, err = template.New("MetricsCollector").Parse(mcs.GoTemplate.RawTemplate)
		} else if mcs.GoTemplate.TemplatePath != "" {
			mtp, err = template.ParseFiles(mcs.GoTemplate.TemplatePath)
		} else {
		}
		if err != nil {
			return nil, err
		}
	}
	if mtp == nil {
		mtp, err = template.ParseFiles("/metricscollector-template/defaultMetricsCollectorTemplate.yaml")
		if err != nil {
			return nil, err
		}
	}
	var b bytes.Buffer
	err = mtp.Execute(&b, tmpValues)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
