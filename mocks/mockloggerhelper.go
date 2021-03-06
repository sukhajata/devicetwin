// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sukhajata/devicetwin.git/pkg/loggerhelper (interfaces: Helper)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	pplogger "github.com/sukhajata/pplogger"
	reflect "reflect"
)

// MockHelper is a mock of Helper interface
type MockHelper struct {
	ctrl     *gomock.Controller
	recorder *MockHelperMockRecorder
}

// MockHelperMockRecorder is the mock recorder for MockHelper
type MockHelperMockRecorder struct {
	mock *MockHelper
}

// NewMockHelper creates a new mock instance
func NewMockHelper(ctrl *gomock.Controller) *MockHelper {
	mock := &MockHelper{ctrl: ctrl}
	mock.recorder = &MockHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHelper) EXPECT() *MockHelperMockRecorder {
	return m.recorder
}

// LogError mocks base method
func (m *MockHelper) LogError(arg0, arg1 string, arg2 pplogger.ErrorMessage_Severity) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogError", arg0, arg1, arg2)
}

// LogError indicates an expected call of LogError
func (mr *MockHelperMockRecorder) LogError(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogError", reflect.TypeOf((*MockHelper)(nil).LogError), arg0, arg1, arg2)
}
