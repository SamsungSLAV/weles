// Code generated by MockGen. DO NOT EDIT.
// Source: git.tizen.org/tools/weles (interfaces: DryadJobManager)

// Package mock is a generated GoMock package.
package mock

import (
	weles "git.tizen.org/tools/weles"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDryadJobManager is a mock of DryadJobManager interface
type MockDryadJobManager struct {
	ctrl     *gomock.Controller
	recorder *MockDryadJobManagerMockRecorder
}

// MockDryadJobManagerMockRecorder is the mock recorder for MockDryadJobManager
type MockDryadJobManagerMockRecorder struct {
	mock *MockDryadJobManager
}

// NewMockDryadJobManager creates a new mock instance
func NewMockDryadJobManager(ctrl *gomock.Controller) *MockDryadJobManager {
	mock := &MockDryadJobManager{ctrl: ctrl}
	mock.recorder = &MockDryadJobManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDryadJobManager) EXPECT() *MockDryadJobManagerMockRecorder {
	return m.recorder
}

// Cancel mocks base method
func (m *MockDryadJobManager) Cancel(arg0 weles.JobID) error {
	ret := m.ctrl.Call(m, "Cancel", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cancel indicates an expected call of Cancel
func (mr *MockDryadJobManagerMockRecorder) Cancel(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockDryadJobManager)(nil).Cancel), arg0)
}

// Create mocks base method
func (m *MockDryadJobManager) Create(arg0 weles.JobID, arg1 weles.Dryad, arg2 weles.Config, arg3 chan<- weles.DryadJobStatusChange) error {
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockDryadJobManagerMockRecorder) Create(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDryadJobManager)(nil).Create), arg0, arg1, arg2, arg3)
}

// List mocks base method
func (m *MockDryadJobManager) List(arg0 *weles.DryadJobFilter) ([]weles.DryadJobInfo, error) {
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]weles.DryadJobInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockDryadJobManagerMockRecorder) List(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDryadJobManager)(nil).List), arg0)
}
