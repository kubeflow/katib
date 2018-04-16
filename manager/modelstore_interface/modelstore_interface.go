package modelstore_interface

import (
	"github.com/kubeflow/hp-tuning/api"
)

type ModelStoreInterface interface {
	StoreStudy(*api.StoreStudyRequest) error
	StoreModel(*api.StoreModelRequest) error
	GetStoredStudies() ([]*api.StudyOverView, error)
	GetStoredModels(*api.GetStoredModelsRequest) ([]*api.ModelInfo, error)
	GetStoredModel(*api.GetStoredModelRequest) (*api.ModelInfo, error)
}
