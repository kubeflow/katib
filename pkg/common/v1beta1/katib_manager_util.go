/*
Copyright 2022 The Kubeflow Authors.

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

package v1beta1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api_pb "github.com/kubeflow/katib/pkg/apis/manager/v1beta1"
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

type katibDBManagerClientAndConn struct {
	Conn                 *grpc.ClientConn
	KatibDBManagerClient api_pb.DBManagerClient
}

// GetDBManagerAddr returns address of Katib DB Manager
func GetDBManagerAddr() string {
	dbManagerNS := consts.DefaultKatibDBManagerServiceNamespace
	dbManagerIP := consts.DefaultKatibDBManagerServiceIP
	dbManagerPort := consts.DefaultKatibDBManagerServicePort

	if len(dbManagerNS) != 0 {
		return dbManagerIP + "." + dbManagerNS + ":" + dbManagerPort
	}

	return dbManagerIP + ":" + dbManagerPort
}

func getKatibDBManagerClientAndConn() (*katibDBManagerClientAndConn, error) {
	addr := GetDBManagerAddr()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	kcc := &katibDBManagerClientAndConn{
		Conn:                 conn,
		KatibDBManagerClient: api_pb.NewDBManagerClient(conn),
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
