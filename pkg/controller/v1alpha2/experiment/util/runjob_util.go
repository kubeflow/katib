/*

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
	"k8s.io/apimachinery/pkg/runtime/schema"
)


func GetSupportdJobList() []schema.GroupVersionKind {
	// TODO: append other supported jobs, such as tfjob, pytorch and so on 
	supportedJobList := []schema.GroupVersionKind {
		schema.GroupVersionKind {
			Group:   "batch",
			Version: "v1",
			Kind:    "Job",
		},
	}
	return supportedJobList
}
