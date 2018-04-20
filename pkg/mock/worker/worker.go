// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeflow/katib/pkg/manager/worker_interface (interfaces: WorkerInterface)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/kubeflow/katib/pkg/api"
	reflect "reflect"
)

// MockWorkerInterface is a mock of WorkerInterface interface
type MockWorkerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerInterfaceMockRecorder
}

// MockWorkerInterfaceMockRecorder is the mock recorder for MockWorkerInterface
type MockWorkerInterfaceMockRecorder struct {
	mock *MockWorkerInterface
}

// NewMockWorkerInterface creates a new mock instance
func NewMockWorkerInterface(ctrl *gomock.Controller) *MockWorkerInterface {
	mock := &MockWorkerInterface{ctrl: ctrl}
	mock.recorder = &MockWorkerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkerInterface) EXPECT() *MockWorkerInterfaceMockRecorder {
	return m.recorder
}

// CheckRunningTrials mocks base method
func (m *MockWorkerInterface) CheckRunningTrials(arg0, arg1 string) error {
	ret := m.ctrl.Call(m, "CheckRunningTrials", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckRunningTrials indicates an expected call of CheckRunningTrials
func (mr *MockWorkerInterfaceMockRecorder) CheckRunningTrials(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRunningTrials", reflect.TypeOf((*MockWorkerInterface)(nil).CheckRunningTrials), arg0, arg1)
}

// CleanWorkers mocks base method
func (m *MockWorkerInterface) CleanWorkers(arg0 string) error {
	ret := m.ctrl.Call(m, "CleanWorkers", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanWorkers indicates an expected call of CleanWorkers
func (mr *MockWorkerInterfaceMockRecorder) CleanWorkers(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanWorkers", reflect.TypeOf((*MockWorkerInterface)(nil).CleanWorkers), arg0)
}

// CompleteTrial mocks base method
func (m *MockWorkerInterface) CompleteTrial(arg0, arg1 string, arg2 bool) error {
	ret := m.ctrl.Call(m, "CompleteTrial", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteTrial indicates an expected call of CompleteTrial
func (mr *MockWorkerInterfaceMockRecorder) CompleteTrial(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteTrial", reflect.TypeOf((*MockWorkerInterface)(nil).CompleteTrial), arg0, arg1, arg2)
}

// GetCompletedTrials mocks base method
func (m *MockWorkerInterface) GetCompletedTrials(arg0 string) []*api.Trial {
	ret := m.ctrl.Call(m, "GetCompletedTrials", arg0)
	ret0, _ := ret[0].([]*api.Trial)
	return ret0
}

// GetCompletedTrials indicates an expected call of GetCompletedTrials
func (mr *MockWorkerInterfaceMockRecorder) GetCompletedTrials(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompletedTrials", reflect.TypeOf((*MockWorkerInterface)(nil).GetCompletedTrials), arg0)
}

// GetRunningTrials mocks base method
func (m *MockWorkerInterface) GetRunningTrials(arg0 string) []*api.Trial {
	ret := m.ctrl.Call(m, "GetRunningTrials", arg0)
	ret0, _ := ret[0].([]*api.Trial)
	return ret0
}

// GetRunningTrials indicates an expected call of GetRunningTrials
func (mr *MockWorkerInterfaceMockRecorder) GetRunningTrials(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunningTrials", reflect.TypeOf((*MockWorkerInterface)(nil).GetRunningTrials), arg0)
}

// GetTrialEvLogs mocks base method
func (m *MockWorkerInterface) GetTrialEvLogs(arg0, arg1 string, arg2 []string, arg3 string) ([]*api.EvaluationLog, error) {
	ret := m.ctrl.Call(m, "GetTrialEvLogs", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*api.EvaluationLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialEvLogs indicates an expected call of GetTrialEvLogs
func (mr *MockWorkerInterfaceMockRecorder) GetTrialEvLogs(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialEvLogs", reflect.TypeOf((*MockWorkerInterface)(nil).GetTrialEvLogs), arg0, arg1, arg2, arg3)
}

// GetTrialObjValue mocks base method
func (m *MockWorkerInterface) GetTrialObjValue(arg0, arg1, arg2 string) (string, error) {
	ret := m.ctrl.Call(m, "GetTrialObjValue", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialObjValue indicates an expected call of GetTrialObjValue
func (mr *MockWorkerInterfaceMockRecorder) GetTrialObjValue(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialObjValue", reflect.TypeOf((*MockWorkerInterface)(nil).GetTrialObjValue), arg0, arg1, arg2)
}

// IsTrialComplete mocks base method
func (m *MockWorkerInterface) IsTrialComplete(arg0, arg1 string) (bool, error) {
	ret := m.ctrl.Call(m, "IsTrialComplete", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsTrialComplete indicates an expected call of IsTrialComplete
func (mr *MockWorkerInterfaceMockRecorder) IsTrialComplete(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTrialComplete", reflect.TypeOf((*MockWorkerInterface)(nil).IsTrialComplete), arg0, arg1)
}

// SpawnWorkers mocks base method
func (m *MockWorkerInterface) SpawnWorkers(arg0 []*api.Trial, arg1 string) error {
	ret := m.ctrl.Call(m, "SpawnWorkers", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SpawnWorkers indicates an expected call of SpawnWorkers
func (mr *MockWorkerInterfaceMockRecorder) SpawnWorkers(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpawnWorkers", reflect.TypeOf((*MockWorkerInterface)(nil).SpawnWorkers), arg0, arg1)
}
