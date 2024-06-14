// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubeflow/katib/pkg/util/v1beta1/katibclient (interfaces: Client)
//
// Generated by this command:
//
//	mockgen -package mock -destination pkg/mock/v1beta1/util/katibclient/katibclient.go github.com/kubeflow/katib/pkg/util/v1beta1/katibclient Client
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	v1beta1 "github.com/kubeflow/katib/pkg/apis/controller/experiments/v1beta1"
	v1beta10 "github.com/kubeflow/katib/pkg/apis/controller/suggestions/v1beta1"
	v1beta11 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	gomock "go.uber.org/mock/gomock"
	v1 "k8s.io/api/core/v1"
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

// CreateRuntimeObject mocks base method.
func (m *MockClient) CreateRuntimeObject(arg0 client.Object) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRuntimeObject", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRuntimeObject indicates an expected call of CreateRuntimeObject.
func (mr *MockClientMockRecorder) CreateRuntimeObject(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRuntimeObject", reflect.TypeOf((*MockClient)(nil).CreateRuntimeObject), arg0)
}

// DeleteRuntimeObject mocks base method.
func (m *MockClient) DeleteRuntimeObject(arg0 client.Object) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRuntimeObject", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRuntimeObject indicates an expected call of DeleteRuntimeObject.
func (mr *MockClientMockRecorder) DeleteRuntimeObject(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRuntimeObject", reflect.TypeOf((*MockClient)(nil).DeleteRuntimeObject), arg0)
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
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConfigMap", varargs...)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigMap indicates an expected call of GetConfigMap.
func (mr *MockClientMockRecorder) GetConfigMap(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigMap", reflect.TypeOf((*MockClient)(nil).GetConfigMap), varargs...)
}

// GetExperiment mocks base method.
func (m *MockClient) GetExperiment(arg0 string, arg1 ...string) (*v1beta1.Experiment, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetExperiment", varargs...)
	ret0, _ := ret[0].(*v1beta1.Experiment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExperiment indicates an expected call of GetExperiment.
func (mr *MockClientMockRecorder) GetExperiment(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExperiment", reflect.TypeOf((*MockClient)(nil).GetExperiment), varargs...)
}

// GetExperimentList mocks base method.
func (m *MockClient) GetExperimentList(arg0 ...string) (*v1beta1.ExperimentList, error) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetExperimentList", varargs...)
	ret0, _ := ret[0].(*v1beta1.ExperimentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExperimentList indicates an expected call of GetExperimentList.
func (mr *MockClientMockRecorder) GetExperimentList(arg0 ...any) *gomock.Call {
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
func (m *MockClient) GetSuggestion(arg0 string, arg1 ...string) (*v1beta10.Suggestion, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSuggestion", varargs...)
	ret0, _ := ret[0].(*v1beta10.Suggestion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSuggestion indicates an expected call of GetSuggestion.
func (mr *MockClientMockRecorder) GetSuggestion(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSuggestion", reflect.TypeOf((*MockClient)(nil).GetSuggestion), varargs...)
}

// GetTrial mocks base method.
func (m *MockClient) GetTrial(arg0 string, arg1 ...string) (*v1beta11.Trial, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrial", varargs...)
	ret0, _ := ret[0].(*v1beta11.Trial)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrial indicates an expected call of GetTrial.
func (mr *MockClientMockRecorder) GetTrial(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrial", reflect.TypeOf((*MockClient)(nil).GetTrial), varargs...)
}

// GetTrialList mocks base method.
func (m *MockClient) GetTrialList(arg0 string, arg1 ...string) (*v1beta11.TrialList, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrialList", varargs...)
	ret0, _ := ret[0].(*v1beta11.TrialList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialList indicates an expected call of GetTrialList.
func (mr *MockClientMockRecorder) GetTrialList(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialList", reflect.TypeOf((*MockClient)(nil).GetTrialList), varargs...)
}

// GetTrialTemplates mocks base method.
func (m *MockClient) GetTrialTemplates(arg0 ...string) (*v1.ConfigMapList, error) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTrialTemplates", varargs...)
	ret0, _ := ret[0].(*v1.ConfigMapList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrialTemplates indicates an expected call of GetTrialTemplates.
func (mr *MockClientMockRecorder) GetTrialTemplates(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrialTemplates", reflect.TypeOf((*MockClient)(nil).GetTrialTemplates), arg0...)
}

// InjectClient mocks base method.
func (m *MockClient) InjectClient(arg0 client.Client) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InjectClient", arg0)
}

// InjectClient indicates an expected call of InjectClient.
func (mr *MockClientMockRecorder) InjectClient(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectClient", reflect.TypeOf((*MockClient)(nil).InjectClient), arg0)
}

// UpdateRuntimeObject mocks base method.
func (m *MockClient) UpdateRuntimeObject(arg0 client.Object) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRuntimeObject", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRuntimeObject indicates an expected call of UpdateRuntimeObject.
func (mr *MockClientMockRecorder) UpdateRuntimeObject(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRuntimeObject", reflect.TypeOf((*MockClient)(nil).UpdateRuntimeObject), arg0)
}
