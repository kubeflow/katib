package modelstore

import (
	"github.com/kubeflow/katib/api"
)

type ModelStore interface {
	SaveStudy(*api.SaveStudyRequest) error
	SaveModel(*api.SaveModelRequest) error
	GetSavedStudies() ([]*api.StudyOverview, error)
	GetSavedModels(*api.GetSavedModelsRequest) ([]*api.ModelInfo, error)
	GetSavedModel(*api.GetSavedModelRequest) (*api.ModelInfo, error)
}
