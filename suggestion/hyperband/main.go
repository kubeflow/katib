package main

import (
	pb "github.com/mlkube/katib/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = "0.0.0.0:6789"
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	size := 1<<31 - 1
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterSuggestionServer(s, NewHyperBandSuggestService())
	reflection.Register(s)
	log.Printf("HyperBand Suggestion Service\n")
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
