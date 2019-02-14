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
	"text/template"

	katibapi "github.com/kubeflow/katib/pkg/api"
	katibv1alpha1 "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"github.com/kubeflow/katib/pkg"
	"github.com/kubeflow/katib/pkg/manager/studyjobclient"

	"k8s.io/apimachinery/pkg/util/uuid"
)

func getTemplateStrFromConfigMap(cNamespace, cName, tPath string) (string, error) {
	sjc, err := studyjobclient.NewStudyjobClient(nil)
	if err != nil {
		return "", err
	}
	return sjc.GetTemplate(cNamespace, cName, tPath)
}

func getTemplateStr(goTemplate *katibv1alpha1.GoTemplate, getDefaultTemplateSpec func()(string,string,string)) (string, error) {
	if goTemplate.RawTemplate != "" {
		return goTemplate.RawTemplate, nil
	} else {
		tName, tNamespace, tPath := getDefaultTemplateSpec()
		if goTemplate.TemplateSpec != nil {
			tName = goTemplate.TemplateSpec.ConfigMapName
			tNamespace = goTemplate.TemplateSpec.ConfigMapNamespace
			tPath = goTemplate.TemplateSpec.TemplatePath
		}
		return getTemplateStrFromConfigMap(tNamespace, tName, tPath)
	}
}

func getDefaultWorkerTemplateSpec() (string, string, string) {
	return getKatibNamespace(), "worker-template", "defaultWorkerTemplate.yaml"
}

func getDefaultMetricsTemplateSpec() (string, string, string) {
	return getKatibNamespace(), "metricscollector-template", "defaultMetricsCollectorTemplate.yaml"
}

func getWorkerTemplateStr(workerSpec *katibv1alpha1.WorkerSpec) (string, error) {
	if workerSpec == nil || workerSpec.GoTemplate == nil {
		return getTemplateStrFromConfigMap(getDefaultWorkerTemplateSpec())
	} else {
		return getTemplateStr(workerSpec.GoTemplate, getDefaultWorkerTemplateSpec)
	}
}

func getWorkerManifest(c katibapi.ManagerClient, studyID string, trial *katibapi.Trial, workerSpec *katibv1alpha1.WorkerSpec, kind string, ns string, dryrun bool) (string, *bytes.Buffer, error) {
	wStr, err := getWorkerTemplateStr(workerSpec)
	if err != nil {
		return "", nil, err
	}
	wtp, err := template.New("Worker").Parse(wStr)
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

func getMetricsCollectorTemplateStr(mcs *katibv1alpha1.MetricsCollectorSpec) (string, error) {
	if mcs == nil || mcs.GoTemplate == nil {
		return getTemplateStrFromConfigMap(getDefaultMetricsTemplateSpec())
	} else {
		return getTemplateStr(mcs.GoTemplate, getDefaultMetricsTemplateSpec)
	}
}

func getMetricsCollectorManifest(studyID string, trialID string, workerID string, workerKind string, namespace string, mcs *katibv1alpha1.MetricsCollectorSpec) (*bytes.Buffer, error) {
	tmpValues := map[string]string{
		"StudyID":    studyID,
		"TrialID":    trialID,
		"WorkerID":   workerID,
		"WorkerKind": workerKind,
		"NameSpace":  namespace,
		"ManagerSerivce": pkg.GetManagerAddr(),
	}
	tStr, err := getMetricsCollectorTemplateStr(mcs)
	if err != nil {
		return nil, err
	}
	mtp, err := template.New("MetricsCollector").Parse(tStr)
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
