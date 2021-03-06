// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sukhajata/devicetwin/internal/dbclient (interfaces: Client)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	types "github.com/sukhajata/devicetwin/internal/types"
	config "github.com/sukhajata/ppconfig"
	reflect "reflect"
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

// DeleteConfig mocks base method
func (m *MockClient) DeleteConfig(arg0 string, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConfig indicates an expected call of DeleteConfig
func (mr *MockClientMockRecorder) DeleteConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConfig", reflect.TypeOf((*MockClient)(nil).DeleteConfig), arg0, arg1)
}

// GetConfigByIndex mocks base method
func (m *MockClient) GetConfigByIndex(arg0 *config.GetConfigByIndexRequest) (*config.ConfigField, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfigByIndex", arg0)
	ret0, _ := ret[0].(*config.ConfigField)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigByIndex indicates an expected call of GetConfigByIndex
func (mr *MockClientMockRecorder) GetConfigByIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigByIndex", reflect.TypeOf((*MockClient)(nil).GetConfigByIndex), arg0)
}

// GetConfigByName mocks base method
func (m *MockClient) GetConfigByName(arg0 string, arg1 types.ConfigFieldDetails, arg2 *config.GetConfigByNameRequest) (*config.ConfigField, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfigByName", arg0, arg1, arg2)
	ret0, _ := ret[0].(*config.ConfigField)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigByName indicates an expected call of GetConfigByName
func (mr *MockClientMockRecorder) GetConfigByName(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigByName", reflect.TypeOf((*MockClient)(nil).GetConfigByName), arg0, arg1, arg2)
}

// GetDLResmin mocks base method
func (m *MockClient) GetDLResmin(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDLResmin", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDLResmin indicates an expected call of GetDLResmin
func (mr *MockClientMockRecorder) GetDLResmin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDLResmin", reflect.TypeOf((*MockClient)(nil).GetDLResmin), arg0)
}

// GetDeviceConfig mocks base method
func (m *MockClient) GetDeviceConfig(arg0 *config.Identifier) (*config.ConfigFields, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeviceConfig", arg0)
	ret0, _ := ret[0].(*config.ConfigFields)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeviceConfig indicates an expected call of GetDeviceConfig
func (mr *MockClientMockRecorder) GetDeviceConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeviceConfig", reflect.TypeOf((*MockClient)(nil).GetDeviceConfig), arg0)
}

// GetFieldDetails mocks base method
func (m *MockClient) GetFieldDetails(arg0, arg1 string) (map[string]types.ConfigFieldDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFieldDetails", arg0, arg1)
	ret0, _ := ret[0].(map[string]types.ConfigFieldDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFieldDetails indicates an expected call of GetFieldDetails
func (mr *MockClientMockRecorder) GetFieldDetails(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFieldDetails", reflect.TypeOf((*MockClient)(nil).GetFieldDetails), arg0, arg1)
}

// GetFieldDetailsByIndex mocks base method
func (m *MockClient) GetFieldDetailsByIndex(arg0 int32, arg1, arg2 string) (types.ConfigFieldDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFieldDetailsByIndex", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.ConfigFieldDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFieldDetailsByIndex indicates an expected call of GetFieldDetailsByIndex
func (mr *MockClientMockRecorder) GetFieldDetailsByIndex(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFieldDetailsByIndex", reflect.TypeOf((*MockClient)(nil).GetFieldDetailsByIndex), arg0, arg1, arg2)
}

// GetFieldDetailsByName mocks base method
func (m *MockClient) GetFieldDetailsByName(arg0, arg1, arg2 string) (types.ConfigFieldDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFieldDetailsByName", arg0, arg1, arg2)
	ret0, _ := ret[0].(types.ConfigFieldDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFieldDetailsByName indicates an expected call of GetFieldDetailsByName
func (mr *MockClientMockRecorder) GetFieldDetailsByName(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFieldDetailsByName", reflect.TypeOf((*MockClient)(nil).GetFieldDetailsByName), arg0, arg1, arg2)
}

// GetInconsistentDevices mocks base method
func (m *MockClient) GetInconsistentDevices() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInconsistentDevices")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInconsistentDevices indicates an expected call of GetInconsistentDevices
func (mr *MockClientMockRecorder) GetInconsistentDevices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInconsistentDevices", reflect.TypeOf((*MockClient)(nil).GetInconsistentDevices))
}

// GetLatestFirmware mocks base method
func (m *MockClient) GetLatestFirmware(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestFirmware", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestFirmware indicates an expected call of GetLatestFirmware
func (mr *MockClientMockRecorder) GetLatestFirmware(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestFirmware", reflect.TypeOf((*MockClient)(nil).GetLatestFirmware), arg0)
}

// GetNextRadioOffset mocks base method
func (m *MockClient) GetNextRadioOffset() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextRadioOffset")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNextRadioOffset indicates an expected call of GetNextRadioOffset
func (mr *MockClientMockRecorder) GetNextRadioOffset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextRadioOffset", reflect.TypeOf((*MockClient)(nil).GetNextRadioOffset))
}

// GetS11ConfigKey mocks base method
func (m *MockClient) GetS11ConfigKey(arg0 string, arg1 int32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetS11ConfigKey", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetS11ConfigKey indicates an expected call of GetS11ConfigKey
func (mr *MockClientMockRecorder) GetS11ConfigKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetS11ConfigKey", reflect.TypeOf((*MockClient)(nil).GetS11ConfigKey), arg0, arg1)
}

// UpdateConfigToNewFirmware mocks base method
func (m *MockClient) UpdateConfigToNewFirmware(arg0 string, arg1 int, arg2 map[string]types.ConfigFieldDetails) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateConfigToNewFirmware", arg0, arg1, arg2)
}

// UpdateConfigToNewFirmware indicates an expected call of UpdateConfigToNewFirmware
func (mr *MockClientMockRecorder) UpdateConfigToNewFirmware(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfigToNewFirmware", reflect.TypeOf((*MockClient)(nil).UpdateConfigToNewFirmware), arg0, arg1, arg2)
}

// UpdateDbDesired mocks base method
func (m *MockClient) UpdateDbDesired(arg0 *config.SetDesiredRequest, arg1 types.ConfigFieldDetails) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDbDesired", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDbDesired indicates an expected call of UpdateDbDesired
func (mr *MockClientMockRecorder) UpdateDbDesired(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDbDesired", reflect.TypeOf((*MockClient)(nil).UpdateDbDesired), arg0, arg1)
}

// UpdateDbReported mocks base method
func (m *MockClient) UpdateDbReported(arg0 *config.UpdateReportedRequest, arg1 types.ConfigFieldDetails) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDbReported", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDbReported indicates an expected call of UpdateDbReported
func (mr *MockClientMockRecorder) UpdateDbReported(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDbReported", reflect.TypeOf((*MockClient)(nil).UpdateDbReported), arg0, arg1)
}

// UpdateFirmwareAllDevices mocks base method
func (m *MockClient) UpdateFirmwareAllDevices() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFirmwareAllDevices")
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFirmwareAllDevices indicates an expected call of UpdateFirmwareAllDevices
func (mr *MockClientMockRecorder) UpdateFirmwareAllDevices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFirmwareAllDevices", reflect.TypeOf((*MockClient)(nil).UpdateFirmwareAllDevices))
}
