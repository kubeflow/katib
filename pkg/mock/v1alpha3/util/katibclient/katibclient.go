// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient (interfaces: Client)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	v1alpha30 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	reflect "reflect"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateExperiment mocks base method
func (m *MockClient) CreateExperiment(arg0 *v1alpha3.Experiment, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateExperiment", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateExperiment indicates an expected call of CreateExperiment
func (mr *MockClientMockRecorder) CreateExperiment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExperiment", reflect.TypeOf((*MockClient)(nil).CreateExperiment), varargs...)
}

// DeleteExperiment mocks base method
func (m *MockClient) DeleteExperiment(arg0 *v1alpha3.Experiment, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteExperiment", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExperiment indicates an expected call of DeleteExperiment
func (mr *MockClientMockRecorder) DeleteExperiment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExperiment", reflect.TypeOf((*MockClient)(nil).DeleteExperiment), varargs...)
}

// GetConfigMap mocks base method
func (m *MockClient) GetConfigMap(arg0 string, arg1 ...string) (map[string]string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConfigMap", varargs...)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigMap indicates an expected call of GetConfigMap
func (mr *MockClientMockRecorder) GetConfigMap(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigMap", reflect.TypeOf((*MockClient)(nil).GetConfigMap), varargs...)
}

// GetExperiment mocks base method
func (m *MockClient) GetExperiment(arg0 string, arg1 ...string) (*v1alpha3.Experiment, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetExperiment", varargs...)
	ret0, _ := ret[0].(*v1alpha3.Experiment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExperiment indicates an expected call of GetExperiment
func (mr *MockClientMockRecorder) GetExperiment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExperiment", reflect.TypeOf((*MockClient)(nil).GetExperiment), varargs...)
}

// GetExperimentList mocks base method
func (m *MockClient) GetExperimentList(arg0 ...string) (*v1alpha3.ExperimentList, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetExperimentList", varargs...)
	ret0, _ := ret[0].(*v1alpha3.ExperimentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExperimentList indicates an expected call of GetExperimentList
func (mr *MockClientMockRecorder) GetExperimentList(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExperimentList", reflect.TypeOf((*MockClient)(nil).GetExperimentList), arg0...)
}

// GetTrialList mocks base method
func (m *MockClient) GetTrialList(arg0 string, arg1 ...string) (*v1alpha30.TrialList, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrialList", varargs...)
	ret0, _ := ret[0].(*v1alpha30.TrialList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialList indicates an expected call of GetTrialList
func (mr *MockClientMockRecorder) GetTrialList(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialList", reflect.TypeOf((*MockClient)(nil).GetTrialList), varargs...)
}

// GetTrialTemplates mocks base method
func (m *MockClient) GetTrialTemplates(arg0 ...string) (map[string]string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrialTemplates", varargs...)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialTemplates indicates an expected call of GetTrialTemplates
func (mr *MockClientMockRecorder) GetTrialTemplates(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialTemplates", reflect.TypeOf((*MockClient)(nil).GetTrialTemplates), arg0...)
}

// InjectClient mocks base method
func (m *MockClient) InjectClient(arg0 client.Client) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InjectClient", arg0)
}

// InjectClient indicates an expected call of InjectClient
func (mr *MockClientMockRecorder) InjectClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectClient", reflect.TypeOf((*MockClient)(nil).InjectClient), arg0)
}

// UpdateTrialTemplates mocks base method
func (m *MockClient) UpdateTrialTemplates(arg0 map[string]string, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateTrialTemplates", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTrialTemplates indicates an expected call of UpdateTrialTemplates
func (mr *MockClientMockRecorder) UpdateTrialTemplates(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrialTemplates", reflect.TypeOf((*MockClient)(nil).UpdateTrialTemplates), varargs...)
}
