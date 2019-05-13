package ui

import "github.com/kubeflow/katib/pkg/util/v1alpha2/katibclient"

const maxMsgSize = 1<<31 - 1

var (
	// namespace      = "default"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
)

type Decoder struct {
	Layers     int            `json:"num_layers"`
	InputSize  []int          `json:"input_size"`
	OutputSize []int          `json:"output_size"`
	Embedding  map[int]*Block `json:"embedding"`
}

type Block struct {
	ID    int    `json:"opt_id"`
	Type  string `json:"opt_type"`
	Param Option `json:"opt_params"`
}

type Option struct {
	FilterNumber string `json:"num_filter"`
	FilterSize   string `json:"filter_size"`
	Stride       string `json:"stride"`
}

type JobView struct {
	Name   string
	Status string
}

type TemplateView struct {
	Name string
	Yaml string
}

type KatibUIHandler struct {
	katibClient *katibclient.KatibClient
}

type TemplateResponse struct {
	TemplateType string
	Data         []TemplateView
}

type NNView struct {
	Name         string
	TrialName    string
	Architecture string
	MetricsName  []string
	MetricsValue []string
}
