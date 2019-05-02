package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "github.com/kubeflow/katib/pkg/api/v1alpha2"
)

var (
	vizierCoreEndpoint = flag.String("echo_endpoint", "vizier-core:6789", "vizier-core endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	// register handlers for the HTTP endpoints
	err := gw.RegisterManagerHandlerFromEndpoint(ctx, mux, *vizierCoreEndpoint, opts)
	if err != nil {
		return err
	}

	// proxy server listens on port 80
	return http.ListenAndServe(":80", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
