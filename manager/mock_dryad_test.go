// Code generated by MockGen. DO NOT EDIT.
// Source: git.tizen.org/tools/weles/manager/dryad (interfaces: SessionProvider,DeviceCommunicationProvider)

// Package manager is a generated GoMock package.
package manager

import (
	dryad "git.tizen.org/tools/weles/manager/dryad"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSessionProvider is a mock of SessionProvider interface
type MockSessionProvider struct {
	ctrl     *gomock.Controller
	recorder *MockSessionProviderMockRecorder
}

// MockSessionProviderMockRecorder is the mock recorder for MockSessionProvider
type MockSessionProviderMockRecorder struct {
	mock *MockSessionProvider
}

// NewMockSessionProvider creates a new mock instance
func NewMockSessionProvider(ctrl *gomock.Controller) *MockSessionProvider {
	mock := &MockSessionProvider{ctrl: ctrl}
	mock.recorder = &MockSessionProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSessionProvider) EXPECT() *MockSessionProviderMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockSessionProvider) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockSessionProviderMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSessionProvider)(nil).Close))
}

// DUT mocks base method
func (m *MockSessionProvider) DUT() error {
	ret := m.ctrl.Call(m, "DUT")
	ret0, _ := ret[0].(error)
	return ret0
}

// DUT indicates an expected call of DUT
func (mr *MockSessionProviderMockRecorder) DUT() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DUT", reflect.TypeOf((*MockSessionProvider)(nil).DUT))
}

// Exec mocks base method
func (m *MockSessionProvider) Exec(arg0 ...string) ([]byte, []byte, error) {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Exec indicates an expected call of Exec
func (mr *MockSessionProviderMockRecorder) Exec(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockSessionProvider)(nil).Exec), arg0...)
}

// PowerTick mocks base method
func (m *MockSessionProvider) PowerTick() error {
	ret := m.ctrl.Call(m, "PowerTick")
	ret0, _ := ret[0].(error)
	return ret0
}

// PowerTick indicates an expected call of PowerTick
func (mr *MockSessionProviderMockRecorder) PowerTick() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PowerTick", reflect.TypeOf((*MockSessionProvider)(nil).PowerTick))
}

// TS mocks base method
func (m *MockSessionProvider) TS() error {
	ret := m.ctrl.Call(m, "TS")
	ret0, _ := ret[0].(error)
	return ret0
}

// TS indicates an expected call of TS
func (mr *MockSessionProviderMockRecorder) TS() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TS", reflect.TypeOf((*MockSessionProvider)(nil).TS))
}

// MockDeviceCommunicationProvider is a mock of DeviceCommunicationProvider interface
type MockDeviceCommunicationProvider struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceCommunicationProviderMockRecorder
}

// MockDeviceCommunicationProviderMockRecorder is the mock recorder for MockDeviceCommunicationProvider
type MockDeviceCommunicationProviderMockRecorder struct {
	mock *MockDeviceCommunicationProvider
}

// NewMockDeviceCommunicationProvider creates a new mock instance
func NewMockDeviceCommunicationProvider(ctrl *gomock.Controller) *MockDeviceCommunicationProvider {
	mock := &MockDeviceCommunicationProvider{ctrl: ctrl}
	mock.recorder = &MockDeviceCommunicationProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDeviceCommunicationProvider) EXPECT() *MockDeviceCommunicationProviderMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockDeviceCommunicationProvider) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockDeviceCommunicationProviderMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDeviceCommunicationProvider)(nil).Close))
}

// CopyFilesFrom mocks base method
func (m *MockDeviceCommunicationProvider) CopyFilesFrom(arg0 []string, arg1 string) error {
	ret := m.ctrl.Call(m, "CopyFilesFrom", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyFilesFrom indicates an expected call of CopyFilesFrom
func (mr *MockDeviceCommunicationProviderMockRecorder) CopyFilesFrom(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyFilesFrom", reflect.TypeOf((*MockDeviceCommunicationProvider)(nil).CopyFilesFrom), arg0, arg1)
}

// CopyFilesTo mocks base method
func (m *MockDeviceCommunicationProvider) CopyFilesTo(arg0 []string, arg1 string) error {
	ret := m.ctrl.Call(m, "CopyFilesTo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyFilesTo indicates an expected call of CopyFilesTo
func (mr *MockDeviceCommunicationProviderMockRecorder) CopyFilesTo(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyFilesTo", reflect.TypeOf((*MockDeviceCommunicationProvider)(nil).CopyFilesTo), arg0, arg1)
}

// Exec mocks base method
func (m *MockDeviceCommunicationProvider) Exec(arg0 ...string) ([]byte, []byte, error) {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Exec indicates an expected call of Exec
func (mr *MockDeviceCommunicationProviderMockRecorder) Exec(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDeviceCommunicationProvider)(nil).Exec), arg0...)
}

// Login mocks base method
func (m *MockDeviceCommunicationProvider) Login(arg0 dryad.Credentials) error {
	ret := m.ctrl.Call(m, "Login", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login
func (mr *MockDeviceCommunicationProviderMockRecorder) Login(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockDeviceCommunicationProvider)(nil).Login), arg0)
}
