// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sukhajata/ppconnection (interfaces: ConnectionServiceClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	ppconnection "github.com/sukhajata/ppconnection"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockConnectionServiceClient is a mock of ConnectionServiceClient interface
type MockConnectionServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionServiceClientMockRecorder
}

// MockConnectionServiceClientMockRecorder is the mock recorder for MockConnectionServiceClient
type MockConnectionServiceClientMockRecorder struct {
	mock *MockConnectionServiceClient
}

// NewMockConnectionServiceClient creates a new mock instance
func NewMockConnectionServiceClient(ctrl *gomock.Controller) *MockConnectionServiceClient {
	mock := &MockConnectionServiceClient{ctrl: ctrl}
	mock.recorder = &MockConnectionServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConnectionServiceClient) EXPECT() *MockConnectionServiceClientMockRecorder {
	return m.recorder
}

// AddSlot mocks base method
func (m *MockConnectionServiceClient) AddSlot(arg0 context.Context, arg1 *ppconnection.AddSlotRequest, arg2 ...grpc.CallOption) (*ppconnection.Identifier, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddSlot", varargs...)
	ret0, _ := ret[0].(*ppconnection.Identifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSlot indicates an expected call of AddSlot
func (mr *MockConnectionServiceClientMockRecorder) AddSlot(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSlot", reflect.TypeOf((*MockConnectionServiceClient)(nil).AddSlot), varargs...)
}

// Cleanup mocks base method
func (m *MockConnectionServiceClient) Cleanup(arg0 context.Context, arg1 *ppconnection.Empty, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Cleanup", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cleanup indicates an expected call of Cleanup
func (mr *MockConnectionServiceClientMockRecorder) Cleanup(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockConnectionServiceClient)(nil).Cleanup), varargs...)
}

// CreateConnection mocks base method
func (m *MockConnectionServiceClient) CreateConnection(arg0 context.Context, arg1 *ppconnection.CreateConnectionRequest, arg2 ...grpc.CallOption) (*ppconnection.Identifier, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateConnection", varargs...)
	ret0, _ := ret[0].(*ppconnection.Identifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateConnection indicates an expected call of CreateConnection
func (mr *MockConnectionServiceClientMockRecorder) CreateConnection(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateConnection", reflect.TypeOf((*MockConnectionServiceClient)(nil).CreateConnection), varargs...)
}

// CreateDevice mocks base method
func (m *MockConnectionServiceClient) CreateDevice(arg0 context.Context, arg1 *ppconnection.Device, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateDevice", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDevice indicates an expected call of CreateDevice
func (mr *MockConnectionServiceClientMockRecorder) CreateDevice(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDevice", reflect.TypeOf((*MockConnectionServiceClient)(nil).CreateDevice), varargs...)
}

// CreateImage mocks base method
func (m *MockConnectionServiceClient) CreateImage(arg0 context.Context, arg1 *ppconnection.ConnectionImage, arg2 ...grpc.CallOption) (*ppconnection.Identifier, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateImage", varargs...)
	ret0, _ := ret[0].(*ppconnection.Identifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateImage indicates an expected call of CreateImage
func (mr *MockConnectionServiceClientMockRecorder) CreateImage(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImage", reflect.TypeOf((*MockConnectionServiceClient)(nil).CreateImage), varargs...)
}

// CreateMultiplePendingConnections mocks base method
func (m *MockConnectionServiceClient) CreateMultiplePendingConnections(arg0 context.Context, arg1 *ppconnection.MultipleConnectionRequest, arg2 ...grpc.CallOption) (*ppconnection.Identifiers, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateMultiplePendingConnections", varargs...)
	ret0, _ := ret[0].(*ppconnection.Identifiers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMultiplePendingConnections indicates an expected call of CreateMultiplePendingConnections
func (mr *MockConnectionServiceClientMockRecorder) CreateMultiplePendingConnections(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMultiplePendingConnections", reflect.TypeOf((*MockConnectionServiceClient)(nil).CreateMultiplePendingConnections), varargs...)
}

// CreatePendingConnection mocks base method
func (m *MockConnectionServiceClient) CreatePendingConnection(arg0 context.Context, arg1 *ppconnection.Connection, arg2 ...grpc.CallOption) (*ppconnection.Identifier, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreatePendingConnection", varargs...)
	ret0, _ := ret[0].(*ppconnection.Identifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePendingConnection indicates an expected call of CreatePendingConnection
func (mr *MockConnectionServiceClientMockRecorder) CreatePendingConnection(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePendingConnection", reflect.TypeOf((*MockConnectionServiceClient)(nil).CreatePendingConnection), varargs...)
}

// DeleteConnection mocks base method
func (m *MockConnectionServiceClient) DeleteConnection(arg0 context.Context, arg1 *ppconnection.Identifier, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteConnection", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteConnection indicates an expected call of DeleteConnection
func (mr *MockConnectionServiceClientMockRecorder) DeleteConnection(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConnection", reflect.TypeOf((*MockConnectionServiceClient)(nil).DeleteConnection), varargs...)
}

// DeleteImage mocks base method
func (m *MockConnectionServiceClient) DeleteImage(arg0 context.Context, arg1 *ppconnection.Identifier, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteImage", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteImage indicates an expected call of DeleteImage
func (mr *MockConnectionServiceClientMockRecorder) DeleteImage(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockConnectionServiceClient)(nil).DeleteImage), varargs...)
}

// GetAddress mocks base method
func (m *MockConnectionServiceClient) GetAddress(arg0 context.Context, arg1 *ppconnection.GetAddressRequest, arg2 ...grpc.CallOption) (*ppconnection.Location, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAddress", varargs...)
	ret0, _ := ret[0].(*ppconnection.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAddress indicates an expected call of GetAddress
func (mr *MockConnectionServiceClientMockRecorder) GetAddress(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetAddress), varargs...)
}

// GetConnection mocks base method
func (m *MockConnectionServiceClient) GetConnection(arg0 context.Context, arg1 *ppconnection.Identifier, arg2 ...grpc.CallOption) (*ppconnection.Connection, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConnection", varargs...)
	ret0, _ := ret[0].(*ppconnection.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnection indicates an expected call of GetConnection
func (mr *MockConnectionServiceClientMockRecorder) GetConnection(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnection", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetConnection), varargs...)
}

// GetConnections mocks base method
func (m *MockConnectionServiceClient) GetConnections(arg0 context.Context, arg1 *ppconnection.GetConnectionsRequest, arg2 ...grpc.CallOption) (*ppconnection.Connections, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConnections", varargs...)
	ret0, _ := ret[0].(*ppconnection.Connections)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnections indicates an expected call of GetConnections
func (mr *MockConnectionServiceClientMockRecorder) GetConnections(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnections", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetConnections), varargs...)
}

// GetConnectionsByTransformer mocks base method
func (m *MockConnectionServiceClient) GetConnectionsByTransformer(arg0 context.Context, arg1 *ppconnection.GetConnectionsByTransformerRequest, arg2 ...grpc.CallOption) (*ppconnection.Connections, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConnectionsByTransformer", varargs...)
	ret0, _ := ret[0].(*ppconnection.Connections)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionsByTransformer indicates an expected call of GetConnectionsByTransformer
func (mr *MockConnectionServiceClientMockRecorder) GetConnectionsByTransformer(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionsByTransformer", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetConnectionsByTransformer), varargs...)
}

// GetConnectionsForIDNumber mocks base method
func (m *MockConnectionServiceClient) GetConnectionsForIDNumber(arg0 context.Context, arg1 *ppconnection.Identifier, arg2 ...grpc.CallOption) (*ppconnection.Connections, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConnectionsForIDNumber", varargs...)
	ret0, _ := ret[0].(*ppconnection.Connections)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionsForIDNumber indicates an expected call of GetConnectionsForIDNumber
func (mr *MockConnectionServiceClientMockRecorder) GetConnectionsForIDNumber(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionsForIDNumber", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetConnectionsForIDNumber), varargs...)
}

// GetDeviceBySerialNumber mocks base method
func (m *MockConnectionServiceClient) GetDeviceBySerialNumber(arg0 context.Context, arg1 *ppconnection.GetDeviceBySerialNumberRequest, arg2 ...grpc.CallOption) (*ppconnection.Device, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetDeviceBySerialNumber", varargs...)
	ret0, _ := ret[0].(*ppconnection.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceBySerialNumber indicates an expected call of GetDeviceBySerialNumber
func (mr *MockConnectionServiceClientMockRecorder) GetDeviceBySerialNumber(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceBySerialNumber", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetDeviceBySerialNumber), varargs...)
}

// GetImage mocks base method
func (m *MockConnectionServiceClient) GetImage(arg0 context.Context, arg1 *ppconnection.Identifier, arg2 ...grpc.CallOption) (*ppconnection.ConnectionImage, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetImage", varargs...)
	ret0, _ := ret[0].(*ppconnection.ConnectionImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImage indicates an expected call of GetImage
func (mr *MockConnectionServiceClientMockRecorder) GetImage(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImage", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetImage), varargs...)
}

// GetLiveConnections mocks base method
func (m *MockConnectionServiceClient) GetLiveConnections(arg0 context.Context, arg1 *ppconnection.GetConnectionsRequest, arg2 ...grpc.CallOption) (*ppconnection.Connections, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLiveConnections", varargs...)
	ret0, _ := ret[0].(*ppconnection.Connections)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLiveConnections indicates an expected call of GetLiveConnections
func (mr *MockConnectionServiceClientMockRecorder) GetLiveConnections(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLiveConnections", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetLiveConnections), varargs...)
}

// GetTransformers mocks base method
func (m *MockConnectionServiceClient) GetTransformers(arg0 context.Context, arg1 *ppconnection.Empty, arg2 ...grpc.CallOption) (*ppconnection.TransformerList, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTransformers", varargs...)
	ret0, _ := ret[0].(*ppconnection.TransformerList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransformers indicates an expected call of GetTransformers
func (mr *MockConnectionServiceClientMockRecorder) GetTransformers(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransformers", reflect.TypeOf((*MockConnectionServiceClient)(nil).GetTransformers), varargs...)
}

// UpdateConnection mocks base method
func (m *MockConnectionServiceClient) UpdateConnection(arg0 context.Context, arg1 *ppconnection.UpdateConnectionRequest, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateConnection", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateConnection indicates an expected call of UpdateConnection
func (mr *MockConnectionServiceClientMockRecorder) UpdateConnection(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConnection", reflect.TypeOf((*MockConnectionServiceClient)(nil).UpdateConnection), varargs...)
}

// UpdateConnectionState mocks base method
func (m *MockConnectionServiceClient) UpdateConnectionState(arg0 context.Context, arg1 *ppconnection.UpdateConnectionStateRequest, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateConnectionState", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateConnectionState indicates an expected call of UpdateConnectionState
func (mr *MockConnectionServiceClientMockRecorder) UpdateConnectionState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConnectionState", reflect.TypeOf((*MockConnectionServiceClient)(nil).UpdateConnectionState), varargs...)
}

// UpdateIdentityTable mocks base method
func (m *MockConnectionServiceClient) UpdateIdentityTable(arg0 context.Context, arg1 *ppconnection.Identifier, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateIdentityTable", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateIdentityTable indicates an expected call of UpdateIdentityTable
func (mr *MockConnectionServiceClientMockRecorder) UpdateIdentityTable(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIdentityTable", reflect.TypeOf((*MockConnectionServiceClient)(nil).UpdateIdentityTable), varargs...)
}

// UpdateJob mocks base method
func (m *MockConnectionServiceClient) UpdateJob(arg0 context.Context, arg1 *ppconnection.UpdateJobRequest, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateJob", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateJob indicates an expected call of UpdateJob
func (mr *MockConnectionServiceClientMockRecorder) UpdateJob(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateJob", reflect.TypeOf((*MockConnectionServiceClient)(nil).UpdateJob), varargs...)
}

// UpdateLines mocks base method
func (m *MockConnectionServiceClient) UpdateLines(arg0 context.Context, arg1 *ppconnection.UpdateLinesRequest, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateLines", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLines indicates an expected call of UpdateLines
func (mr *MockConnectionServiceClientMockRecorder) UpdateLines(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLines", reflect.TypeOf((*MockConnectionServiceClient)(nil).UpdateLines), varargs...)
}

// UpdateMount mocks base method
func (m *MockConnectionServiceClient) UpdateMount(arg0 context.Context, arg1 *ppconnection.UpdateMountRequest, arg2 ...grpc.CallOption) (*ppconnection.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateMount", varargs...)
	ret0, _ := ret[0].(*ppconnection.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMount indicates an expected call of UpdateMount
func (mr *MockConnectionServiceClientMockRecorder) UpdateMount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMount", reflect.TypeOf((*MockConnectionServiceClient)(nil).UpdateMount), varargs...)
}
