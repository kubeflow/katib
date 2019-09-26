/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha3

import (
	"context"
	"os"

	"google.golang.org/grpc"

	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

const (
	KatibManagerServiceIPEnvName        = "KATIB_MANAGER_PORT_6789_TCP_ADDR"
	KatibManagerServicePortEnvName      = "KATIB_MANAGER_PORT_6789_TCP_PORT"
	KatibManagerServiceNamespaceEnvName = "KATIB_MANAGER_NAMESPACE"
	KatibManagerService                 = "katib-manager"
	KatibManagerPort                    = "6789"
	ManagerAddr                         = KatibManagerService + ":" + KatibManagerPort
)

type katibManagerClientAndConn struct {
	Conn               *grpc.ClientConn
	KatibManagerClient api_pb.ManagerClient
}

func GetManagerAddr() string {
	ns := consts.DefaultKatibNamespace
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

func getKatibManagerClientAndConn() (*katibManagerClientAndConn, error) {
	addr := GetManagerAddr()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	kcc := &katibManagerClientAndConn{
		Conn:               conn,
		KatibManagerClient: api_pb.NewManagerClient(conn),
	}
	return kcc, nil
}

func closeKatibManagerConnection(kcc *katibManagerClientAndConn) {
	kcc.Conn.Close()
}

func GetObservationLog(request *api_pb.GetObservationLogRequest) (*api_pb.GetObservationLogReply, error) {
	ctx := context.Background()
	kcc, err := getKatibManagerClientAndConn()
	if err != nil {
		return nil, err
	}
	defer closeKatibManagerConnection(kcc)
	kc := kcc.KatibManagerClient
	return kc.GetObservationLog(ctx, request)
}
