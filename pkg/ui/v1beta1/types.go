package v1beta1

import (
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
	"github.com/kubeflow/katib/pkg/util/v1beta1/katibclient"
)

const maxMsgSize = 1<<31 - 1

var (
	// namespace      = "default"
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"

	TrialTemplateLabel = map[string]string{
		consts.LabelTrialTemplateConfigMapName: consts.LabelTrialTemplateConfigMapValue}
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
	Name      string
	Status    string
	Namespace string
}

type TrialTemplatesDataView struct {
	ConfigMapNamespace string
	ConfigMaps         []ConfigMap
}

type TrialTemplatesResponse struct {
	Data []TrialTemplatesDataView
}

type ConfigMap struct {
	ConfigMapName string
	Templates     []Template
}

type Template struct {
	Path string
	Yaml string
}

type KatibUIHandler struct {
	katibClient   katibclient.Client
	dbManagerAddr string
}

type NNView struct {
	Name         string
	TrialName    string
	Architecture string
	MetricsName  []string
	MetricsValue []string
}

type JobType string

const (
	JobTypeHP        = "HP"
	JobTypeNAS       = "NAS"
	ActionTypeAdd    = "add"
	ActionTypeEdit   = "edit"
	ActionTypeDelete = "delete"
)
