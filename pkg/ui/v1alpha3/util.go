package v1alpha3

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	gographviz "github.com/awalterschulze/gographviz"
)

func (k *KatibUIHandler) getExperimentList(namespace string, typ JobType) ([]JobView, error) {
	jobs := make([]JobView, 0)

	el, err := k.katibClient.GetExperimentList(namespace)
	if err != nil {
		log.Printf("GetExperimentList failed: %v", err)
		return nil, err
	}
	for _, experiment := range el.Items {
		if (typ == JobTypeNAS && experiment.Spec.NasConfig != nil) ||
			(typ == JobTypeHP && experiment.Spec.NasConfig == nil) {
			experimentLastCondition, err := experiment.GetLastConditionType()
			if err != nil {
				log.Printf("GetLastConditionType failed: %v", err)
				return nil, err
			}
			jobs = append(jobs, JobView{
				Name:      experiment.Name,
				Namespace: experiment.Namespace,
				Status:    string(experimentLastCondition),
			})
		}
	}
	return jobs, nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "text/html; charset=utf-8")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	(*w).Header().Set("Access-Control-Expose-Headers", "Access-Control-*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func getTemplatesView(templates map[string]string) []TemplateView {
	templatesView := make([]TemplateView, 0)

	for key := range templates {
		templatesView = append(templatesView, TemplateView{Name: key, Yaml: templates[key]})
	}
	return templatesView
}

func (k *KatibUIHandler) updateTemplates(newTemplate map[string]interface{}, isDelete bool) (TemplateResponse, error) {
	var currentTemplates map[string]string
	var err error

	currentTemplates, err = k.katibClient.GetTrialTemplates()
	if err != nil {
		return TemplateResponse{}, errors.New("GetTrialTemplates failed: " + err.Error())
	}

	if isDelete {
		delete(currentTemplates, newTemplate["name"].(string))
	} else {
		currentTemplates[newTemplate["name"].(string)] = newTemplate["yaml"].(string)
	}

	err = k.katibClient.UpdateTrialTemplates(currentTemplates)
	if err != nil {
		return TemplateResponse{}, errors.New("UpdateTrialTemplates failed: " + err.Error())
	}

	TemplateResponse := TemplateResponse{
		Data:         getTemplatesView(currentTemplates),
		TemplateType: newTemplate["kind"].(string),
	}
	return TemplateResponse, nil
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
