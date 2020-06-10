package util

import (
	"bytes"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
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
