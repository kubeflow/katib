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

package util

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/kubeflow/katib/pkg"
	katibv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	"github.com/kubeflow/katib/pkg/util/v1alpha2/katibclient"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getMetricsCollectorManifest(experimentID string, trialID string, jobKind string, namespace string, mcs *katibv1alpha2.MetricsCollectorSpec) (*bytes.Buffer, error) {
	var mtp *template.Template = nil
	var err error
	tmpValues := map[string]string{
		"ExperimentID":   experimentID,
		"TrialID":        trialID,
		"JobKind":        jobKind,
		"NameSpace":      namespace,
		"ManagerService": pkg.GetManagerAddr(),
	}
	if mcs != nil && mcs.GoTemplate.RawTemplate != "" {
		mtp, err = template.New("MetricsCollector").Parse(mcs.GoTemplate.RawTemplate)
	} else {
		mctp := "defaultMetricsCollectorTemplate.yaml"
		if mcs != nil && mcs.GoTemplate.TemplateSpec != nil {
			mctp = mcs.GoTemplate.TemplateSpec.TemplatePath
		}
		kc, err := katibclient.NewClient(client.Options{})
		if err != nil {
			return nil, err
		}
		mtl, err := kc.GetMetricsCollectorTemplates()
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
