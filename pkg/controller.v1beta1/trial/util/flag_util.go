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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GvkListFlag is the custom flag to parse GroupVersionKind list for trial resources.
type GvkListFlag []schema.GroupVersionKind

// Set is the method to convert gvk to string value
func (flag *GvkListFlag) String() string {
	gvkStrings := []string{}
	for _, x := range []schema.GroupVersionKind(*flag) {
		gvkStrings = append(gvkStrings, x.String())
	}
	return strings.Join(gvkStrings, ",")
}

// Set is the method to set gvk from string flag value
func (flag *GvkListFlag) Set(value string) error {
	gvk, _ := schema.ParseKindArg(value)
	if gvk == nil {
		return fmt.Errorf("Invalid GroupVersionKind: %v", value)
	}
	*flag = append(*flag, *gvk)
	return nil
}
