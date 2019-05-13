package v1alpha2

import (
	"os"
	"context"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"google.golang.org/grpc"
)

const (
	KatibManagerServiceIPEnvName        = "KATIB_MANAGER_PORT_6789_TCP_ADDR"
	KatibManagerServicePortEnvName      = "KATIB_MANAGER_PORT_6789_TCP_PORT"
	KatibManagerServiceNamespaceEnvName = "KATIB_MANAGER_NAMESPACE"
	KatibManagerService                 = "katib-manager"
	KatibManagerPort                    = "6789"
	ManagerAddr                   = KatibManagerService + ":" + KatibManagerPort
)

type katibClientAndConnection struct {
	Conn *grpc.ClientConn
	KatibClient api_pb.ManagerClient
}

func GetManagerAddr() string {
	ns := os.Getenv(experimentsv1alpha2.DefaultKatibNamespaceEnvName)
	if len(ns) == 0 {
		addr := os.Getenv(KatibManagerServiceIPEnvName)
		port := os.Getenv(KatibManagerServicePortEnvName)
		if len(addr) > 0 && len(port) > 0 {
			return addr + ":" + port
		} else {
			return ManagerAddr
		}
	} else {
		return KatibManagerService + "." + ns + ":" + KatibManagerPort
	}
}

func getKatibClientAndConnection() (*katibClientAndConnection, error) {
	addr := GetManagerAddr()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	kc := &katibClientAndConnection {
		Conn: conn,
		KatibClient: api_pb.NewManagerClient(conn),
	}
	return kc, nil
}

func RegisterExperiment(request *api_pb.RegisterExperimentRequest) (*api_pb.RegisterExperimentReply, error) {
	ctx := context.Background()
	kcc, err := getKatibClientAndConnection()
	if err != nil {
		return nil, err
	}
	defer kcc.Conn.Close()
	kc := kcc.KatibClient
	return kc.RegisterExperiment(ctx, request)
}
