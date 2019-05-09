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
	"bytes"
	"context"
	"fmt"
	"text/template"

	common "github.com/kubeflow/katib/pkg/common/v1alpha1"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	katibapi "github.com/kubeflow/katib/pkg/api/v1alpha1"
	"github.com/kubeflow/katib/pkg/manager/v1alpha1/studyjobclient"

	"k8s.io/apimachinery/pkg/util/uuid"
)

func getWorkerManifest(c katibapi.ManagerClient, studyID string, trial *katibapi.Trial, workerSpec *katibv1alpha1.WorkerSpec, kind string, ns string, dryrun bool) (string, *bytes.Buffer, error) {
	var wtp *template.Template = nil
	var err error
	if workerSpec != nil && workerSpec.GoTemplate.RawTemplate != "" {
		wtp, err = template.New("Worker").Parse(workerSpec.GoTemplate.RawTemplate)
	} else {
		wPath := "defaultWorkerTemplate.yaml"
		if workerSpec != nil && workerSpec.GoTemplate.TemplatePath != "" {
			wPath = workerSpec.GoTemplate.TemplatePath
		}
		sjc, err := studyjobclient.NewStudyjobClient(nil)
		if err != nil {
			return "", nil, err
		}
		wtl, err := sjc.GetWorkerTemplates()
		if err != nil {
			return "", nil, err
		}
		if wt, ok := wtl[wPath]; !ok {
			return "", nil, fmt.Errorf("No worker template name %s", wPath)
		} else {
			wtp, err = template.New("Worker").Parse(wt)
		}
	}
	if err != nil {
		return "", nil, err
	}
	var wid string
	if dryrun {
		wid = string(uuid.NewUUID())
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
		StudyID:   studyID,
		TrialID:   trial.TrialId,
		WorkerID:  wid,
		NameSpace: ns,
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

func getMetricsCollectorManifest(studyID string, trialID string, workerID string, workerKind string, namespace string, mcs *katibv1alpha1.MetricsCollectorSpec) (*bytes.Buffer, error) {
	var mtp *template.Template = nil
	var err error
	tmpValues := map[string]string{
		"StudyID":        studyID,
		"TrialID":        trialID,
		"WorkerID":       workerID,
		"WorkerKind":     workerKind,
		"NameSpace":      namespace,
		"ManagerSerivce": common.GetManagerAddr(),
	}
	if mcs != nil && mcs.GoTemplate.RawTemplate != "" {
		mtp, err = template.New("MetricsCollector").Parse(mcs.GoTemplate.RawTemplate)
	} else {
		mctp := "defaultMetricsCollectorTemplate.yaml"
		if mcs != nil && mcs.GoTemplate.TemplatePath != "" {
			mctp = mcs.GoTemplate.TemplatePath
		}
		sjc, err := studyjobclient.NewStudyjobClient(nil)
		if err != nil {
			return nil, err
		}
		mtl, err := sjc.GetMetricsCollectorTemplates()
		if err != nil {
			return nil, err
		}
		if mt, ok := mtl[mctp]; !ok {
			return nil, fmt.Errorf("No MetricsCollector template name %s", mctp)
		} else {
			mtp, err = template.New("MetricsCollector").Parse(mt)
		}
	}
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	err = mtp.Execute(&b, tmpValues)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
