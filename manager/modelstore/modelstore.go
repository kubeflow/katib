package modelstore

import (
	"github.com/kubeflow/hp-tuning/api"
)

type ModelSave interface {
	SaveStudy(*api.SaveStudyRequest) error
	SaveModel(*api.SaveModelRequest) error
	GetSavedStudies() ([]*api.StudyOverview, error)
	GetSavedModels(*api.GetSavedModelsRequest) ([]*api.ModelInfo, error)
	GetSavedModel(*api.GetSavedModelRequest) (*api.ModelInfo, error)
}
