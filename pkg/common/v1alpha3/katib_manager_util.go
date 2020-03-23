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

	"google.golang.org/grpc"

	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

type katibDBManagerClientAndConn struct {
	Conn                 *grpc.ClientConn
	KatibDBManagerClient api_pb.ManagerClient
}

// GetDBManagerAddr returns address of Katib DB Manager
func GetDBManagerAddr() string {

	if len(consts.DefaultKatibDBManagerNamespace) != 0 {
		return consts.DefaultKatibDBManagerIP + "." + consts.DefaultKatibDBManagerNamespace + ":" + consts.DefaultKatibDBManagerPort
	}

	return consts.DefaultKatibDBManagerIP + ":" + consts.DefaultKatibDBManagerPort
}

func getKatibDBManagerClientAndConn() (*katibDBManagerClientAndConn, error) {
	addr := GetDBManagerAddr()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	kcc := &katibDBManagerClientAndConn{
		Conn:                 conn,
		KatibDBManagerClient: api_pb.NewManagerClient(conn),
	}
	return kcc, nil
}

func closeKatibDBManagerConnection(kcc *katibDBManagerClientAndConn) {
	kcc.Conn.Close()
}

func GetObservationLog(request *api_pb.GetObservationLogRequest) (*api_pb.GetObservationLogReply, error) {
	ctx := context.Background()
	kcc, err := getKatibDBManagerClientAndConn()
	if err != nil {
		return nil, err
	}
	defer closeKatibDBManagerConnection(kcc)
	kc := kcc.KatibDBManagerClient
	return kc.GetObservationLog(ctx, request)
}

func DeleteObservationLog(request *api_pb.DeleteObservationLogRequest) (*api_pb.DeleteObservationLogReply, error) {
	ctx := context.Background()
	kcc, err := getKatibDBManagerClientAndConn()
	if err != nil {
		return nil, err
	}
	defer closeKatibDBManagerConnection(kcc)
	kc := kcc.KatibDBManagerClient
	return kc.DeleteObservationLog(ctx, request)
}
