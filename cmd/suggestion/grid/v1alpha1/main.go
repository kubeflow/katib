package main

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kubeflow/katib/pkg/api/v1alpha1"
	suggestion "github.com/kubeflow/katib/pkg/suggestion/v1alpha1"
	"k8s.io/klog"
)

func main() {
	listener, err := net.Listen("tcp", suggestion.DefaultPort)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}
	size := 1<<31 - 1
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterSuggestionServer(s, suggestion.NewGridSuggestService())
	reflection.Register(s)
	klog.Info("Grid Search Suggestion Service")
	if err = s.Serve(listener); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
