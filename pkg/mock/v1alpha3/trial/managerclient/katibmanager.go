// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeflow/katib/pkg/controller.v1alpha3/trial/managerclient (interfaces: ManagerClient)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	v1alpha30 "github.com/kubeflow/katib/pkg/apis/manager/v1alpha3"
	reflect "reflect"
)

// MockManagerClient is a mock of ManagerClient interface
type MockManagerClient struct {
	ctrl     *gomock.Controller
	recorder *MockManagerClientMockRecorder
}

// MockManagerClientMockRecorder is the mock recorder for MockManagerClient
type MockManagerClientMockRecorder struct {
	mock *MockManagerClient
}

// NewMockManagerClient creates a new mock instance
func NewMockManagerClient(ctrl *gomock.Controller) *MockManagerClient {
	mock := &MockManagerClient{ctrl: ctrl}
	mock.recorder = &MockManagerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManagerClient) EXPECT() *MockManagerClientMockRecorder {
	return m.recorder
}

// GetTrialObservationLog mocks base method
func (m *MockManagerClient) GetTrialObservationLog(arg0 *v1alpha3.Trial) (*v1alpha30.GetObservationLogReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrialObservationLog", arg0)
	ret0, _ := ret[0].(*v1alpha30.GetObservationLogReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialObservationLog indicates an expected call of GetTrialObservationLog
func (mr *MockManagerClientMockRecorder) GetTrialObservationLog(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialObservationLog", reflect.TypeOf((*MockManagerClient)(nil).GetTrialObservationLog), arg0)
}
