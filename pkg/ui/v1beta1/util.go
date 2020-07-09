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
	apiv1 "k8s.io/api/core/v1"
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
		newExperiment := ExperimentView{
			Name:      experiment.Name,
			Namespace: experiment.Namespace,
			Status:    string(experimentLastCondition),
		}

		if experiment.Spec.NasConfig == nil {
			newExperiment.Type = ExperimentTypeHP
		} else {
			newExperiment.Type = ExperimentTypeNAS
		}

		experiments = append(experiments, newExperiment)
	}

	return experiments, nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "text/html; charset=utf-8")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	(*w).Header().Set("Access-Control-Expose-Headers", "Access-Control-*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func (k *KatibUIHandler) getTrialTemplatesViewList() ([]TrialTemplatesDataView, error) {
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

		if len(trialTemplatesConfigMapList.Items) != 0 {
			trialTemplatesDataView = append(trialTemplatesDataView, getTrialTemplatesView(trialTemplatesConfigMapList))
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

func getTrialTemplatesView(templatesConfigMapList *apiv1.ConfigMapList) TrialTemplatesDataView {

	trialTemplatesDataView := TrialTemplatesDataView{
		ConfigMapNamespace: templatesConfigMapList.Items[0].ObjectMeta.Namespace,
		ConfigMaps:         []ConfigMap{},
	}
	for _, configMap := range templatesConfigMapList.Items {
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

		trialTemplatesDataView.ConfigMaps = append(trialTemplatesDataView.ConfigMaps, newConfigMap)
	}

	return trialTemplatesDataView
}

func (k *KatibUIHandler) updateTrialTemplates(
	updatedConfigMapNamespace,
	updatedConfigMapName,
	configMapPath,
	updatedConfigMapPath,
	updatedTemplateYaml,
	actionType string) ([]TrialTemplatesDataView, error) {

	templates, err := k.katibClient.GetConfigMap(updatedConfigMapName, updatedConfigMapNamespace)
	if err != nil {
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
	} else {
		err = k.katibClient.UpdateRuntimeObject(templatesConfigMap)
		if err != nil {
			log.Printf("UpdateRuntimeObject failed: %v", err)
			return nil, err
		}
	}

	newTemplates, err := k.getTrialTemplatesViewList()
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

		Beforehand substite all ' to " and wrap the string in `
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
	graph.AddNode("G", "0", map[string]string{"label": strconv.Quote("Input")})
	var i int
	for i = 0; i < len(architectureInt); i++ {
		graph.AddNode("G", strconv.Itoa(i+1), map[string]string{"label": getNodeString(decoderParsed.Embedding[architectureInt[i][0]])})
		graph.AddEdge(strconv.Itoa(i), strconv.Itoa(i+1), true, nil)
		for j := 1; j < i+1; j++ {
			if architectureInt[i][j] == 1 {
				graph.AddEdge(strconv.Itoa(j-1), strconv.Itoa(i+1), true, nil)
			}
		}
	}
	graph.AddNode("G", strconv.Itoa(i+1), map[string]string{"label": strconv.Quote("GlobalAvgPool")})
	graph.AddEdge(strconv.Itoa(i), strconv.Itoa(i+1), true, nil)
	graph.AddNode("G", strconv.Itoa(i+2), map[string]string{"label": strconv.Quote("FullConnect\nSoftmax")})
	graph.AddEdge(strconv.Itoa(i+1), strconv.Itoa(i+2), true, nil)
	graph.AddNode("G", strconv.Itoa(i+3), map[string]string{"label": strconv.Quote("Output")})
	graph.AddEdge(strconv.Itoa(i+2), strconv.Itoa(i+3), true, nil)
	s := graph.String()
	return s
}
