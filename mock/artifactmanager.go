// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SamsungSLAV/weles (interfaces: ArtifactManager)

// Package mock is a generated GoMock package.
package mock

import (
	weles "github.com/SamsungSLAV/weles"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockArtifactManager is a mock of ArtifactManager interface
type MockArtifactManager struct {
	ctrl     *gomock.Controller
	recorder *MockArtifactManagerMockRecorder
}

// MockArtifactManagerMockRecorder is the mock recorder for MockArtifactManager
type MockArtifactManagerMockRecorder struct {
	mock *MockArtifactManager
}

// NewMockArtifactManager creates a new mock instance
func NewMockArtifactManager(ctrl *gomock.Controller) *MockArtifactManager {
	mock := &MockArtifactManager{ctrl: ctrl}
	mock.recorder = &MockArtifactManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockArtifactManager) EXPECT() *MockArtifactManagerMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockArtifactManager) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockArtifactManagerMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockArtifactManager)(nil).Close))
}

// CreateArtifact mocks base method
func (m *MockArtifactManager) CreateArtifact(arg0 weles.ArtifactDescription) (weles.ArtifactPath, error) {
	ret := m.ctrl.Call(m, "CreateArtifact", arg0)
	ret0, _ := ret[0].(weles.ArtifactPath)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateArtifact indicates an expected call of CreateArtifact
func (mr *MockArtifactManagerMockRecorder) CreateArtifact(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateArtifact", reflect.TypeOf((*MockArtifactManager)(nil).CreateArtifact), arg0)
}

// GetArtifactInfo mocks base method
func (m *MockArtifactManager) GetArtifactInfo(arg0 weles.ArtifactPath) (weles.ArtifactInfo, error) {
	ret := m.ctrl.Call(m, "GetArtifactInfo", arg0)
	ret0, _ := ret[0].(weles.ArtifactInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArtifactInfo indicates an expected call of GetArtifactInfo
func (mr *MockArtifactManagerMockRecorder) GetArtifactInfo(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArtifactInfo", reflect.TypeOf((*MockArtifactManager)(nil).GetArtifactInfo), arg0)
}

// ListArtifact mocks base method
func (m *MockArtifactManager) ListArtifact(arg0 weles.ArtifactFilter, arg1 weles.ArtifactSorter, arg2 weles.ArtifactPaginator) ([]weles.ArtifactInfo, weles.ListInfo, error) {
	ret := m.ctrl.Call(m, "ListArtifact", arg0, arg1, arg2)
	ret0, _ := ret[0].([]weles.ArtifactInfo)
	ret1, _ := ret[1].(weles.ListInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListArtifact indicates an expected call of ListArtifact
func (mr *MockArtifactManagerMockRecorder) ListArtifact(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListArtifact", reflect.TypeOf((*MockArtifactManager)(nil).ListArtifact), arg0, arg1, arg2)
}

// PushArtifact mocks base method
func (m *MockArtifactManager) PushArtifact(arg0 weles.ArtifactDescription, arg1 chan weles.ArtifactStatusChange) (weles.ArtifactPath, error) {
	ret := m.ctrl.Call(m, "PushArtifact", arg0, arg1)
	ret0, _ := ret[0].(weles.ArtifactPath)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushArtifact indicates an expected call of PushArtifact
func (mr *MockArtifactManagerMockRecorder) PushArtifact(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushArtifact", reflect.TypeOf((*MockArtifactManager)(nil).PushArtifact), arg0, arg1)
}
