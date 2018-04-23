package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kubeflow/katib/pkg/api"
	"github.com/kubeflow/katib/pkg/suggestion"
)

func main() {
	listener, err := net.Listen("tcp", suggestion.DefaultPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	size := 1<<31 - 1
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterSuggestionServer(s, suggestion.NewRandomSuggestService())
	reflection.Register(s)
	log.Println("Random Suggestion Service")
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
