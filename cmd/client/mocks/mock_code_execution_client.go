// Code generated by MockGen. DO NOT EDIT.
// Source: lokesh-katari/code-realm/cmd/client/codeExecutionpb (interfaces: CodeExecutionServiceClient)

// Package mock_codeExecutionpb is a generated GoMock package.
package mock_codeExecutionpb

import (
	context "context"
	codeExecutionpb "lokesh-katari/code-realm/cmd/client/codeExecutionpb"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockCodeExecutionServiceClient is a mock of CodeExecutionServiceClient interface.
type MockCodeExecutionServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCodeExecutionServiceClientMockRecorder
}

// MockCodeExecutionServiceClientMockRecorder is the mock recorder for MockCodeExecutionServiceClient.
type MockCodeExecutionServiceClientMockRecorder struct {
	mock *MockCodeExecutionServiceClient
}

// NewMockCodeExecutionServiceClient creates a new mock instance.
func NewMockCodeExecutionServiceClient(ctrl *gomock.Controller) *MockCodeExecutionServiceClient {
	mock := &MockCodeExecutionServiceClient{ctrl: ctrl}
	mock.recorder = &MockCodeExecutionServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodeExecutionServiceClient) EXPECT() *MockCodeExecutionServiceClientMockRecorder {
	return m.recorder
}

// ExecuteCode mocks base method.
func (m *MockCodeExecutionServiceClient) ExecuteCode(arg0 context.Context, arg1 *codeExecutionpb.ExecuteCodeRequest, arg2 ...grpc.CallOption) (*codeExecutionpb.ExecuteCodeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecuteCode", varargs...)
	ret0, _ := ret[0].(*codeExecutionpb.ExecuteCodeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteCode indicates an expected call of ExecuteCode.
func (mr *MockCodeExecutionServiceClientMockRecorder) ExecuteCode(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteCode", reflect.TypeOf((*MockCodeExecutionServiceClient)(nil).ExecuteCode), varargs...)
}