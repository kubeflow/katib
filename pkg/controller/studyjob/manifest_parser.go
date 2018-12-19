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
	"log"
	"text/template"

	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"github.com/kubeflow/katib/pkg/manager/studyjobclient"

	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func getWorkerKind(workerSpec *katibv1alpha1.WorkerSpec) (string, error) {
	var typeChecker interface{}
	BUFSIZE := 1024
	_, m, err := getWorkerManifest(
		nil,
		"validation",
		&katibapi.Trial{
			TrialId:      "validation",
			ParameterSet: []*katibapi.Parameter{},
		},
		workerSpec,
		"",
		"",
		true,
	)
	if err != nil {
		return "", err
	}
	if err := k8syaml.NewYAMLOrJSONDecoder(m, BUFSIZE).Decode(&typeChecker); err != nil {
		log.Printf("Yaml decode validation error %v", err)
		return "", err
	}
	tcMap, ok := typeChecker.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	wkind, ok := tcMap["kind"]
	if !ok {
		return "", fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	wkindS, ok := wkind.(string)
	if !ok {
		return "", fmt.Errorf("Cannot get kind of worker %v", typeChecker)
	}
	for _, kind := range ValidWorkerKindList {
		if kind == wkindS {
			return wkindS, nil
		}
	}
	return "", fmt.Errorf("Invalid kind of worker %v", typeChecker)
}

func getWorkerManifest(c katibapi.ManagerClient, studyID string, trial *katibapi.Trial, workerSpec *katibv1alpha1.WorkerSpec, kind string, ns string, dryrun bool) (string, *bytes.Buffer, error) {
	var wtp *template.Template = nil
	var err error
	if workerSpec != nil {
		if workerSpec.GoTemplate.RawTemplate != "" {
			wtp, err = template.New("Worker").Parse(workerSpec.GoTemplate.RawTemplate)
		} else if workerSpec.GoTemplate.TemplatePath != "" {
			sjc, err := studyjobclient.NewStudyjobClient(nil)
			if err != nil {
				return "", nil, err
			}
			wtl, err := sjc.GetWorkerTemplates()
			if err != nil {
				return "", nil, err
			}
			if wt, ok := wtl[workerSpec.GoTemplate.TemplatePath]; !ok {
				return "", nil, fmt.Errorf("No tamplate name %s", workerSpec.GoTemplate.TemplatePath)
			} else {
				wtp, err = template.New("Worker").Parse(wt)
			}
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
		"StudyID":    studyID,
		"TrialID":    trialID,
		"WorkerID":   workerID,
		"WorkerKind": workerKind,
		"NameSpace":  namespace,
	}
	mctp := "defaultMetricsCollectorTemplate.yaml"
	if mcs != nil {
		if mcs.GoTemplate.RawTemplate != "" {
			mtp, err = template.New("MetricsCollector").Parse(mcs.GoTemplate.RawTemplate)
		} else if mcs.GoTemplate.TemplatePath != "" {
			mctp = mcs.GoTemplate.TemplatePath
		}
	} else {
		sjc, err := studyjobclient.NewStudyjobClient(nil)
		if err != nil {
			return nil, err
		}
		mtl, err := sjc.GetMetricsCollectorTemplates()
		if err != nil {
			return nil, err
		}
		if mt, ok := mtl[mctp]; !ok {
			return nil, fmt.Errorf("No tamplate name %s", mctp)
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
