package main

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog"

	pb "github.com/kubeflow/katib/pkg/api/v1alpha1"
	earlystopping "github.com/kubeflow/katib/pkg/earlystopping/v1alpha1"
)

func main() {
	listener, err := net.Listen("tcp", earlystopping.DefaultPort)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}
	size := 1<<31 - 1
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterEarlyStoppingServer(s, earlystopping.NewMedianStoppingRule())
	reflection.Register(s)
	klog.Info("Median Stopping Rule EarlyStopping Service")
	if err = s.Serve(listener); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
