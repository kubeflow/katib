package v1alpha3

import (
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
	"github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient"
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

type TrialTemplatesView struct {
	Namespace      string
	ConfigMapsList []ConfigMapsList
}

type TrialTemplatesResponse struct {
	Data []TrialTemplatesView
}

type ConfigMapsList struct {
	ConfigMapName string
	TemplatesList []TemplatesList
}

type TemplatesList struct {
	Name string
	Yaml string
}

type KatibUIHandler struct {
	katibClient katibclient.Client
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
