/*
Copyright 2021 The Kubeflow Authors.

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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/kubeflow/katib/pkg/apis/v1beta1"
	"k8s.io/klog"
	"k8s.io/kube-openapi/pkg/common"
)

// Generate OpenAPI spec definitions for Katib Resource
func main() {
	if len(os.Args) <= 2 {
		klog.Fatal("Supply Swagger version and Katib Version")
	}
	version := os.Args[1]
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	refCallback := func(name string) spec.Ref {
		return spec.MustCreateRef("#/definitions/" + common.EscapeJsonPointer(swaggify(name)))
	}

	katibVersion := os.Args[2]
	oAPIDefs := make(map[string]common.OpenAPIDefinition)
	if katibVersion == "v1beta1" {
		oAPIDefs = v1beta1.GetOpenAPIDefinitions(refCallback)
	} else {
		klog.Fatalf("Katib version %v is not supported", katibVersion)
	}

	defs := spec.Definitions{}
	for defName, val := range oAPIDefs {
		defs[swaggify(defName)] = val.Schema
	}
	swagger := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:     "2.0",
			Definitions: defs,
			Paths:       &spec.Paths{Paths: map[string]spec.PathItem{}},
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "Katib",
					Description: "Swagger description for Katib",
					Version:     version,
				},
			},
		},
	}
	jsonBytes, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		klog.Fatal(err.Error())
	}
	fmt.Println(string(jsonBytes))
}

func swaggify(name string) string {
	name = strings.Replace(name, "github.com/kubeflow/katib/pkg/apis/controller/common/", "", -1)
	name = strings.Replace(name, "github.com/kubeflow/katib/pkg/apis/controller/experiments/", "", -1)
	name = strings.Replace(name, "github.com/kubeflow/katib/pkg/apis/controller/suggestions", "", -1)
	name = strings.Replace(name, "github.com/kubeflow/katib/pkg/apis/controller/trials", "", -1)
	name = strings.Replace(name, "k8s.io/api/core/", "", -1)
	name = strings.Replace(name, "k8s.io/apimachinery/pkg/apis/meta/", "", -1)
	name = strings.Replace(name, "/", ".", -1)
	return name
}
