/*
Copyright 2018 The Kubernetes Authors.

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

package parse

import (
	"fmt"
	"reflect"
	"testing"

	"k8s.io/gengo/types"
)

func TestParseScaleParams(t *testing.T) {
	testCases := []struct {
		name     string
		tag      string
		expected map[string]string
		parseErr error
	}{
		{
			name: "test ok",
			tag:  "+kubebuilder:subresource:scale:specpath=.spec.replica,statuspath=.status.replica,selectorpath=.spec.Label",
			expected: map[string]string{
				specReplicasPath:   ".spec.replica",
				statusReplicasPath: ".status.replica",
				labelSelectorPath:  ".spec.Label",
			},
			parseErr: nil,
		},
		{
			name: "test ok without selectorpath",
			tag:  "+kubebuilder:subresource:scale:specpath=.spec.replica,statuspath=.status.replica",
			expected: map[string]string{
				specReplicasPath:   ".spec.replica",
				statusReplicasPath: ".status.replica",
			},
			parseErr: nil,
		},
		{
			name: "test ok selectorpath has empty value",
			tag:  "+kubebuilder:subresource:scale:specpath=.spec.replica,statuspath=.status.replica,selectorpath=",
			expected: map[string]string{
				specReplicasPath:   ".spec.replica",
				statusReplicasPath: ".status.replica",
				labelSelectorPath:  "",
			},
			parseErr: nil,
		},
		{
			name:     "test no jsonpath",
			tag:      "+kubebuilder:subresource:scale",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
		{
			name:     "test no specpath",
			tag:      "+kubebuilder:subresource:scale:statuspath=.status.replica,selectorpath=.jsonpath",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
		{
			name:     "test no statuspath",
			tag:      "+kubebuilder:subresource:scale:specpath=.spec.replica,selectorpath=.jsonpath",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
		{
			name:     "test statuspath is empty string",
			tag:      "+kubebuilder:subresource:scale:statuspath=,selectorpath=.jsonpath",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
		{
			name:     "test scale jsonpath has incorrect separator",
			tag:      "+kubebuilder:subresource:scale,specpath=.spec.replica,statuspath=.jsonpath,selectorpath=.jsonpath",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
		{
			name:     "test scale jsonpath has extra separator",
			tag:      "+kubebuilder:subresource:scale:specpath=.spec.replica,statuspath=.status.replicas,selectorpath=.jsonpath,",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
		{
			name:     "test scale jsonpath has incorrect separator in-between key value pairs",
			tag:      "+kubebuilder:subresource:scale:specpath=.spec.replica;statuspath=.jsonpath;selectorpath=.jsonpath",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
		{
			name:     "test unsupported key value pairs",
			tag:      "+kubebuilder:subresource:scale:name=test,specpath=.spec.replica,statuspath=.status.replicas,selectorpath=.jsonpath",
			expected: nil,
			parseErr: fmt.Errorf(jsonPathError),
		},
	}

	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		r := &types.Type{}
		r.CommentLines = []string{tc.tag}
		res, err := parseScaleParams(r)
		if !reflect.DeepEqual(err, tc.parseErr) {
			t.Errorf("test [%s] failed. error is (%v),\n but expected (%v)", tc.name, err, tc.parseErr)
		}
		if !reflect.DeepEqual(res, tc.expected) {
			t.Errorf("test [%s] failed. result is (%v),\n but expected (%v)", tc.name, res, tc.expected)
		}
	}
}
