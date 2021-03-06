// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sukhajata/pplogger (interfaces: LoggerServiceClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	pplogger "github.com/sukhajata/pplogger"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockLoggerServiceClient is a mock of LoggerServiceClient interface
type MockLoggerServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerServiceClientMockRecorder
}

// MockLoggerServiceClientMockRecorder is the mock recorder for MockLoggerServiceClient
type MockLoggerServiceClientMockRecorder struct {
	mock *MockLoggerServiceClient
}

// NewMockLoggerServiceClient creates a new mock instance
func NewMockLoggerServiceClient(ctrl *gomock.Controller) *MockLoggerServiceClient {
	mock := &MockLoggerServiceClient{ctrl: ctrl}
	mock.recorder = &MockLoggerServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLoggerServiceClient) EXPECT() *MockLoggerServiceClientMockRecorder {
	return m.recorder
}

// LogDeviceEvent mocks base method
func (m *MockLoggerServiceClient) LogDeviceEvent(arg0 context.Context, arg1 *pplogger.DeviceLogMessage, arg2 ...grpc.CallOption) (*pplogger.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LogDeviceEvent", varargs...)
	ret0, _ := ret[0].(*pplogger.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogDeviceEvent indicates an expected call of LogDeviceEvent
func (mr *MockLoggerServiceClientMockRecorder) LogDeviceEvent(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDeviceEvent", reflect.TypeOf((*MockLoggerServiceClient)(nil).LogDeviceEvent), varargs...)
}

// LogError mocks base method
func (m *MockLoggerServiceClient) LogError(arg0 context.Context, arg1 *pplogger.ErrorMessage, arg2 ...grpc.CallOption) (*pplogger.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LogError", varargs...)
	ret0, _ := ret[0].(*pplogger.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogError indicates an expected call of LogError
func (mr *MockLoggerServiceClientMockRecorder) LogError(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogError", reflect.TypeOf((*MockLoggerServiceClient)(nil).LogError), varargs...)
}

// LogOpAlarm mocks base method
func (m *MockLoggerServiceClient) LogOpAlarm(arg0 context.Context, arg1 *pplogger.OpAlarmMessage, arg2 ...grpc.CallOption) (*pplogger.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LogOpAlarm", varargs...)
	ret0, _ := ret[0].(*pplogger.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogOpAlarm indicates an expected call of LogOpAlarm
func (mr *MockLoggerServiceClientMockRecorder) LogOpAlarm(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogOpAlarm", reflect.TypeOf((*MockLoggerServiceClient)(nil).LogOpAlarm), varargs...)
}