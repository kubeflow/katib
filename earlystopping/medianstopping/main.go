package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kubeflow/hp-tuning/api"
	"github.com/kubeflow/hp-tuning/earlystopping"
)

func main() {
	listener, err := net.Listen("tcp", earlystopping.DefaultPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	size := 1<<31 - 1
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterEarlyStoppingServer(s, earlystopping.NewMedianStoppingRule())
	reflection.Register(s)
	log.Printf("Median Stopping Rule EarlyStopping Service\n")
	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
