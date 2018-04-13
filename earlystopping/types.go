package earlystopping

import (
	"context"

	"github.com/kubeflow/hp-tuning/api"
)

const (
	// DefaultPort is the port to serve the earlystopping service.
	DefaultPort = "0.0.0.0:6789"
)

// EarlyStoppingService is the interface for earlystopping service.
type EarlyStoppingService interface {
	ShouldTrialStop(ctx context.Context, in *api.ShouldTrialStopRequest) (*api.ShouldTrialStopReply, error)
	SetEarlyStoppingParameter(ctx context.Context, in *api.SetEarlyStoppingParameterRequest) (*api.SetEarlyStoppingParameterReply, error)
}
