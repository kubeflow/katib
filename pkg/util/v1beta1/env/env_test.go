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

package env

import (
	"fmt"
	"os"
	"testing"
)

func TestGetEnvWithDefault(t *testing.T) {
	expected := "FAKE"
	key := "TEST"
	v := GetEnvOrDefault(key, expected)
	if v != expected {
		t.Errorf("Expected %s, got %s", expected, v)
	}
	expected = "FAKE1"
	os.Setenv(key, expected)
	v = GetEnvOrDefault(key, "")
	if v != expected {
		t.Errorf("Expected %s, got %s", expected, v)
	}
}

func TestGetBoolEnvWithDefault(t *testing.T) {
	expected := false
	key := "TEST"
	v := GetBoolEnvOrDefault(key, expected)
	if v != expected {
		t.Errorf("Expected %t, got %t", expected, v)
	}

	expected = true
	envValue := fmt.Sprintf("%t", expected)
	os.Setenv(key, envValue)
	v = GetBoolEnvOrDefault(key, false)
	if v != expected {
		t.Errorf("Expected %t, got %t", expected, v)
	}
}
