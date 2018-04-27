package main

import (
	"github.com/kubeflow/katib/pkg/api"
)

type StudyData struct {
	StudyConf *api.StudyConfig
	Models    []*api.ModelInfo
}
