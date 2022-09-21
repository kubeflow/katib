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

package util

import (
	"bytes"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	bufferSize = 1024
)

var (
	logUnstructured = logf.Log.WithName("util-unstructured")
)

// ConvertUnstructuredToString returns string from Unstructured value
func ConvertUnstructuredToString(in *unstructured.Unstructured) (string, error) {
	inByte, err := in.MarshalJSON()
	if err != nil {
		logUnstructured.Error(err, "MarshalJSON failed")
		return "", err
	}

	return string(inByte), nil
}

// ConvertStringToUnstructured returns Unstructured from string value
func ConvertStringToUnstructured(in string) (*unstructured.Unstructured, error) {
	inBytes := bytes.NewBufferString(in)
	out := &unstructured.Unstructured{}

	err := k8syaml.NewYAMLOrJSONDecoder(inBytes, bufferSize).Decode(out)
	if err != nil {
		logUnstructured.Error(err, "Decode Unstructured to String failed")
		return nil, err
	}

	return out, nil
}

// ConvertObjectToUnstructured returns Unstructured from Kubernetes Object value
func ConvertObjectToUnstructured(in interface{}) (*unstructured.Unstructured, error) {
	out := &unstructured.Unstructured{}
	var err error

	out.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&in)
	if err != nil {
		logUnstructured.Error(err, "Convert Object to Unstructured failed")
		return nil, err
	}

	return out, nil
}
