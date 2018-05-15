package suggestion

import (
	"context"

	"github.com/kubeflow/katib/pkg/api"
)

const (
	// DefaultPort is the port to serve the suggestion service.
	DefaultPort = "0.0.0.0:6789"
	manager     = "vizier-core:6789"
)

// SuggestService is the interface for suggestion service.
type SuggestService interface {
	GetSuggestions(ctx context.Context, in *api.GetSuggestionsRequest) (*api.GetSuggestionsReply, error)
}
