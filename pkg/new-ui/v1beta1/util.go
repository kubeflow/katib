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

package v1beta1

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"

	gographviz "github.com/awalterschulze/gographviz"
	trialv1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *KatibUIHandler) getExperiments(namespace []string) ([]ExperimentView, error) {
	experiments := []ExperimentView{}

	el, err := k.katibClient.GetExperimentList(namespace...)
	if err != nil {
		log.Printf("GetExperimentList failed: %v", err)
		return nil, err
	}
	for _, experiment := range el.Items {
		experimentLastCondition, err := experiment.GetLastConditionType()
		if err != nil {
			log.Printf("GetLastConditionType failed: %v", err)
			return nil, err
		}
		tp := ExperimentTypeHP
		if experiment.Spec.NasConfig != nil {
			tp = ExperimentTypeNAS
		}

		newExperiment := ExperimentView{
			experiment.Name,
			experiment.Namespace,
			tp,
			string(experimentLastCondition),
			experiment.Status,
		}

		experiments = append(experiments, newExperiment)
	}

	return experiments, nil
}

func havePipelineUID(trials []trialv1beta1.Trial) bool {
	for _, t := range trials {
		_, ok := t.GetAnnotations()[kfpRunIDAnnotation]
		if ok {
			return true
		}
	}

	return false
}

func (k *KatibUIHandler) getTrialTemplatesViewList(r *http.Request) ([]TrialTemplatesDataView, error) {
	trialTemplatesDataView := make([]TrialTemplatesDataView, 0)

	// Get all available namespaces
	namespaces, err := k.getAvailableNamespaces()
	if err != nil {
		log.Printf("GetAvailableNamespaces failed: %v", err)
		return nil, err
	}

	// Get Trial Template ConfigMap for each namespace
	for _, ns := range namespaces {
		trialTemplatesConfigMapList, err := k.katibClient.GetTrialTemplates(ns)
		if err != nil {
			log.Printf("GetTrialTemplates failed: %v", err)
			return nil, err
		}

		// Iterate over the trialTemplatesConfigMapList from all namespaces and filter out the
		// configmaps that belong to namespaces other than kubeflow and the ones that the user has
		// access privileges.
		var newTrialTemplatesConfigMapList []apiv1.ConfigMap
		for _, cmap := range trialTemplatesConfigMapList.Items {
			if ns == consts.DefaultKatibNamespace {
				// Add the configmaps from kubeflow namespace (no user can have access to kubeflow ns, so we hardcode it)
				newTrialTemplatesConfigMapList = append(newTrialTemplatesConfigMapList, cmap)
			} else {
				// for all other namespaces check authorization rbac
				configmapName := cmap.ObjectMeta.Name
				user, err := IsAuthorized(consts.ActionTypeGet, ns, apiv1.ResourceConfigMaps.String(), "", configmapName, apiv1.SchemeGroupVersion, k.katibClient.GetClient(), r)
				if err != nil {
					log.Printf("The user: %s is not authorized to view configMap: %s in namespace: %s \n", user, configmapName, ns)
					return nil, err
				} else {
					log.Printf("The user: %s is authorized to view configMap: %s in namespace: %s", user, configmapName, ns)
					newTrialTemplatesConfigMapList = append(newTrialTemplatesConfigMapList, cmap)
				}
			}

		}

		if len(trialTemplatesConfigMapList.Items) != 0 {
			newTrialTemplatesView := getTrialTemplatesView(newTrialTemplatesConfigMapList)
			// ConfigMap with templates must exists in namespace
			if len(newTrialTemplatesView.ConfigMaps) > 0 {
				trialTemplatesDataView = append(trialTemplatesDataView, newTrialTemplatesView)
			}
		}
	}
	return trialTemplatesDataView, nil
}

func (k *KatibUIHandler) getAvailableNamespaces() ([]string, error) {
	var namespaces []string

	namespaceList, err := k.katibClient.GetNamespaceList()
	if err != nil {
		namespaces = append(namespaces, consts.DefaultKatibNamespace)
		return namespaces, nil
	}
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, ns.ObjectMeta.Name)
	}

	return namespaces, nil
}

func getTrialTemplatesView(templatesConfigMapList []apiv1.ConfigMap) TrialTemplatesDataView {

	trialTemplatesDataView := TrialTemplatesDataView{
		ConfigMapNamespace: templatesConfigMapList[0].ObjectMeta.Namespace,
		ConfigMaps:         []ConfigMap{},
	}
	for _, configMap := range templatesConfigMapList {
		newConfigMap := ConfigMap{
			ConfigMapName: configMap.ObjectMeta.Name,
			Templates:     []Template{},
		}
		for key := range configMap.Data {
			newTemplate := Template{
				Path: key,
				Yaml: configMap.Data[key],
			}
			newConfigMap.Templates = append(newConfigMap.Templates, newTemplate)
		}

		// Sort Trial Templates by Path
		sort.SliceStable(newConfigMap.Templates, func(i, j int) bool {
			return newConfigMap.Templates[i].Path <= newConfigMap.Templates[j].Path
		})

		// Templates with data must exists in ConfigMap
		if len(newConfigMap.Templates) > 0 {
			trialTemplatesDataView.ConfigMaps = append(trialTemplatesDataView.ConfigMaps, newConfigMap)
		}

	}
	return trialTemplatesDataView
}

func (k *KatibUIHandler) updateTrialTemplates(
	updatedConfigMapNamespace,
	updatedConfigMapName,
	configMapPath,
	updatedConfigMapPath,
	updatedTemplateYaml,
	actionType string,
	r *http.Request) ([]TrialTemplatesDataView, error) {

	templates, err := k.katibClient.GetConfigMap(updatedConfigMapName, updatedConfigMapNamespace)
	if err != nil && !(errors.IsNotFound(err) && actionType == ActionTypeAdd) {
		log.Printf("GetConfigMap failed: %v", err)
		return nil, err
	}

	if actionType == ActionTypeAdd {
		if len(templates) == 0 {
			templates = make(map[string]string)
			templates[updatedConfigMapPath] = updatedTemplateYaml
		} else {
			templates[updatedConfigMapPath] = updatedTemplateYaml
		}
	} else if actionType == ActionTypeEdit {
		delete(templates, configMapPath)
		templates[updatedConfigMapPath] = updatedTemplateYaml
	} else if actionType == ActionTypeDelete {
		delete(templates, updatedConfigMapPath)
	}

	templatesConfigMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      updatedConfigMapName,
			Namespace: updatedConfigMapNamespace,
			Labels:    TrialTemplateLabel,
		},
		Data: templates,
	}

	// If templates is empty delete Trial template configMap
	if len(templates) == 0 {
		err = k.katibClient.DeleteRuntimeObject(templatesConfigMap)
		if err != nil {
			log.Printf("DeleteRuntimeObject failed: %v", err)
			return nil, err
		}
		// If len(templates) == 1 and adding template, we must create new ConfigMap
	} else if len(templates) == 1 && actionType == ActionTypeAdd {
		err = k.katibClient.CreateRuntimeObject(templatesConfigMap)
		if err != nil {
			log.Printf("CreateRuntimeObject failed: %v", err)
			return nil, err
		}
		// Otherwise updating configMap
	} else {
		err = k.katibClient.UpdateRuntimeObject(templatesConfigMap)
		if err != nil {
			log.Printf("UpdateRuntimeObject failed: %v", err)
			return nil, err
		}
	}

	newTemplates, err := k.getTrialTemplatesViewList(r)
	if err != nil {
		log.Printf("getTrialTemplatesViewList: %v", err)
		return nil, err
	}

	return newTemplates, nil

}
func getNodeString(block *Block) string {
	var nodeString string
	switch block.Type {
	case "convolution":
		nodeString += block.Param.FilterSize + "x" + block.Param.FilterSize
		nodeString += " conv\n"
		nodeString += block.Param.FilterSize + " channels"
	case "separable_convolution":
		nodeString += block.Param.FilterSize + "x" + block.Param.FilterSize
		nodeString += " sep_conv\n"
		nodeString += block.Param.FilterSize + " channels"
	case "depthwise_convolution":
		nodeString += block.Param.FilterSize + "x" + block.Param.FilterSize
		nodeString += " depth_conv\n"
	case "reduction":
		// TODO: Need to be fixed
		nodeString += "3x3 max_pooling"
	}
	return strconv.Quote(nodeString)
}

func generateNNImage(architecture string, decoder string) string {

	var architectureInt [][]int

	if err := json.Unmarshal([]byte(architecture), &architectureInt); err != nil {
		panic(err)
	}
	/*
		Always has num_layers, input_size, output_size and embeding
		Embeding is a map: int to Parameter
		Parameter has id, type, Option

		Beforehand substitute all ' to " and wrap the string in `
	*/

	replacedDecoder := strings.Replace(decoder, `'`, `"`, -1)
	var decoderParsed Decoder

	err := json.Unmarshal([]byte(replacedDecoder), &decoderParsed)
	if err != nil {
		panic(err)
	}

	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	if err = graph.AddNode("G", "0", map[string]string{"label": strconv.Quote("Input")}); err != nil {
		panic(err)
	}
	var i int
	for i = 0; i < len(architectureInt); i++ {
		if err = graph.AddNode("G", strconv.Itoa(i+1), map[string]string{"label": getNodeString(decoderParsed.Embedding[architectureInt[i][0]])}); err != nil {
			panic(err)
		}
		if err = graph.AddEdge(strconv.Itoa(i), strconv.Itoa(i+1), true, nil); err != nil {
			panic(err)
		}
		for j := 1; j < i+1; j++ {
			if architectureInt[i][j] == 1 {
				if err = graph.AddEdge(strconv.Itoa(j-1), strconv.Itoa(i+1), true, nil); err != nil {
					panic(err)
				}
			}
		}
	}
	if err = graph.AddNode("G", strconv.Itoa(i+1), map[string]string{"label": strconv.Quote("GlobalAvgPool")}); err != nil {
		panic(err)
	}
	if err = graph.AddEdge(strconv.Itoa(i), strconv.Itoa(i+1), true, nil); err != nil {
		panic(err)
	}
	if err = graph.AddNode("G", strconv.Itoa(i+2), map[string]string{"label": strconv.Quote("FullConnect\nSoftmax")}); err != nil {
		panic(err)
	}
	if err = graph.AddEdge(strconv.Itoa(i+1), strconv.Itoa(i+2), true, nil); err != nil {
		panic(err)
	}
	if err = graph.AddNode("G", strconv.Itoa(i+3), map[string]string{"label": strconv.Quote("Output")}); err != nil {
		panic(err)
	}
	if err = graph.AddEdge(strconv.Itoa(i+2), strconv.Itoa(i+3), true, nil); err != nil {
		panic(err)
	}
	s := graph.String()
	return s
}
