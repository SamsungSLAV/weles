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

// Create mocks base method
func (m *MockArtifactManager) Create(arg0 weles.ArtifactDescription) (weles.ArtifactPath, error) {
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(weles.ArtifactPath)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockArtifactManagerMockRecorder) Create(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArtifactManager)(nil).Create), arg0)
}

// Download mocks base method
func (m *MockArtifactManager) Download(arg0 weles.ArtifactDescription, arg1 chan weles.ArtifactStatusChange) (weles.ArtifactPath, error) {
	ret := m.ctrl.Call(m, "Download", arg0, arg1)
	ret0, _ := ret[0].(weles.ArtifactPath)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download
func (mr *MockArtifactManagerMockRecorder) Download(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockArtifactManager)(nil).Download), arg0, arg1)
}

// ListArtifact mocks base method
func (m *MockArtifactManager) ListArtifact(arg0 weles.ArtifactFilter, arg1 weles.ArtifactSorter, arg2 weles.ArtifactPagination) ([]weles.ArtifactInfo, weles.ListInfo, error) {
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
