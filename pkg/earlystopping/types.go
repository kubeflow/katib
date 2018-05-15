package earlystopping

import (
	"context"

	"github.com/kubeflow/katib/pkg/api"
)

const (
	// DefaultPort is the port to serve the earlystopping service.
	DefaultPort = "0.0.0.0:6789"
)

// EarlyStoppingService is the interface for earlystopping service.
type EarlyStoppingService interface {
	GetEarlyStoppingParameters(ctx context.Context, in *api.GetEarlyStoppingParametersRequest) (api.GetEarlyStoppingParametersReply, error)
}
