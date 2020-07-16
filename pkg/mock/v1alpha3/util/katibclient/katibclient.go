// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeflow/katib/pkg/util/v1alpha3/katibclient (interfaces: Client)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha3 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1alpha3"
	v1alpha30 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1alpha3"
	v1alpha31 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1alpha3"
	v1 "k8s.io/api/core/v1"
	reflect "reflect"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateExperiment mocks base method.
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

// CreateExperiment indicates an expected call of CreateExperiment.
func (mr *MockClientMockRecorder) CreateExperiment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExperiment", reflect.TypeOf((*MockClient)(nil).CreateExperiment), varargs...)
}

// DeleteExperiment mocks base method.
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

// DeleteExperiment indicates an expected call of DeleteExperiment.
func (mr *MockClientMockRecorder) DeleteExperiment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExperiment", reflect.TypeOf((*MockClient)(nil).DeleteExperiment), varargs...)
}

// GetClient mocks base method.
func (m *MockClient) GetClient() client.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClient")
	ret0, _ := ret[0].(client.Client)
	return ret0
}

// GetClient indicates an expected call of GetClient.
func (mr *MockClientMockRecorder) GetClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClient", reflect.TypeOf((*MockClient)(nil).GetClient))
}

// GetConfigMap mocks base method.
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

// GetConfigMap indicates an expected call of GetConfigMap.
func (mr *MockClientMockRecorder) GetConfigMap(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigMap", reflect.TypeOf((*MockClient)(nil).GetConfigMap), varargs...)
}

// GetExperiment mocks base method.
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

// GetExperiment indicates an expected call of GetExperiment.
func (mr *MockClientMockRecorder) GetExperiment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExperiment", reflect.TypeOf((*MockClient)(nil).GetExperiment), varargs...)
}

// GetExperimentList mocks base method.
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

// GetExperimentList indicates an expected call of GetExperimentList.
func (mr *MockClientMockRecorder) GetExperimentList(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExperimentList", reflect.TypeOf((*MockClient)(nil).GetExperimentList), arg0...)
}

// GetNamespaceList mocks base method.
func (m *MockClient) GetNamespaceList() (*v1.NamespaceList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespaceList")
	ret0, _ := ret[0].(*v1.NamespaceList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespaceList indicates an expected call of GetNamespaceList.
func (mr *MockClientMockRecorder) GetNamespaceList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespaceList", reflect.TypeOf((*MockClient)(nil).GetNamespaceList))
}

// GetSuggestion mocks base method.
func (m *MockClient) GetSuggestion(arg0 string, arg1 ...string) (*v1alpha30.Suggestion, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSuggestion", varargs...)
	ret0, _ := ret[0].(*v1alpha30.Suggestion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSuggestion indicates an expected call of GetSuggestion.
func (mr *MockClientMockRecorder) GetSuggestion(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSuggestion", reflect.TypeOf((*MockClient)(nil).GetSuggestion), varargs...)
}

// GetTrial mocks base method.
func (m *MockClient) GetTrial(arg0 string, arg1 ...string) (*v1alpha31.Trial, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrial", varargs...)
	ret0, _ := ret[0].(*v1alpha31.Trial)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrial indicates an expected call of GetTrial.
func (mr *MockClientMockRecorder) GetTrial(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrial", reflect.TypeOf((*MockClient)(nil).GetTrial), varargs...)
}

// GetTrialList mocks base method.
func (m *MockClient) GetTrialList(arg0 string, arg1 ...string) (*v1alpha31.TrialList, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrialList", varargs...)
	ret0, _ := ret[0].(*v1alpha31.TrialList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialList indicates an expected call of GetTrialList.
func (mr *MockClientMockRecorder) GetTrialList(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialList", reflect.TypeOf((*MockClient)(nil).GetTrialList), varargs...)
}

// GetTrialTemplates mocks base method.
func (m *MockClient) GetTrialTemplates(arg0 ...string) (*v1.ConfigMapList, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrialTemplates", varargs...)
	ret0, _ := ret[0].(*v1.ConfigMapList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialTemplates indicates an expected call of GetTrialTemplates.
func (mr *MockClientMockRecorder) GetTrialTemplates(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialTemplates", reflect.TypeOf((*MockClient)(nil).GetTrialTemplates), arg0...)
}

// InjectClient mocks base method.
func (m *MockClient) InjectClient(arg0 client.Client) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InjectClient", arg0)
}

// InjectClient indicates an expected call of InjectClient.
func (mr *MockClientMockRecorder) InjectClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectClient", reflect.TypeOf((*MockClient)(nil).InjectClient), arg0)
}

// UpdateConfigMap mocks base method.
func (m *MockClient) UpdateConfigMap(arg0 *v1.ConfigMap) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfigMap", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfigMap indicates an expected call of UpdateConfigMap.
func (mr *MockClientMockRecorder) UpdateConfigMap(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfigMap", reflect.TypeOf((*MockClient)(nil).UpdateConfigMap), arg0)
}

// UpdateExperiment mocks base method.
func (m *MockClient) UpdateExperiment(arg0 *v1alpha3.Experiment, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateExperiment", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExperiment indicates an expected call of UpdateExperiment.
func (mr *MockClientMockRecorder) UpdateExperiment(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExperiment", reflect.TypeOf((*MockClient)(nil).UpdateExperiment), varargs...)
}
