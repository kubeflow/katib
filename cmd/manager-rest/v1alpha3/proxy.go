package main

import (
	"flag"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/klog"

	gw "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
)

var (
	katibManagerEndpoint = flag.String("echo_endpoint", "katib-manager:6789", "katib-manager endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	// register handlers for the HTTP endpoints
	err := gw.RegisterManagerHandlerFromEndpoint(ctx, mux, *katibManagerEndpoint, opts)
	if err != nil {
		return err
	}

	// proxy server listens on port 80
	return http.ListenAndServe(":80", mux)
}

func main() {
	flag.Parse()
	defer klog.Flush()

	if err := run(); err != nil {
		klog.Fatal(err)
	}
}
