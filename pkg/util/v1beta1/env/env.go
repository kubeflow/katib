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
	"k8s.io/klog"
	"os"
	"strconv"
)

func GetEnvOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetBoolEnvOrDefault(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		parsedValue, err := strconv.ParseBool(value)
		if err != nil {
			klog.Fatalf("Failed converting %s env to bool", key)
		}
		return parsedValue
	}
	return fallback
}
