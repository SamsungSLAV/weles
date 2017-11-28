// Code generated by MockGen. DO NOT EDIT.
// Source: git.tizen.org/tools/weles (interfaces: JobManager)

// Package mock is a generated GoMock package.
package mock

import (
	weles "git.tizen.org/tools/weles"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockJobManager is a mock of JobManager interface
type MockJobManager struct {
	ctrl     *gomock.Controller
	recorder *MockJobManagerMockRecorder
}

// MockJobManagerMockRecorder is the mock recorder for MockJobManager
type MockJobManagerMockRecorder struct {
	mock *MockJobManager
}

// NewMockJobManager creates a new mock instance
func NewMockJobManager(ctrl *gomock.Controller) *MockJobManager {
	mock := &MockJobManager{ctrl: ctrl}
	mock.recorder = &MockJobManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJobManager) EXPECT() *MockJobManagerMockRecorder {
	return m.recorder
}

// CancelJob mocks base method
func (m *MockJobManager) CancelJob(arg0 weles.JobID) error {
	ret := m.ctrl.Call(m, "CancelJob", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelJob indicates an expected call of CancelJob
func (mr *MockJobManagerMockRecorder) CancelJob(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelJob", reflect.TypeOf((*MockJobManager)(nil).CancelJob), arg0)
}

// CreateJob mocks base method
func (m *MockJobManager) CreateJob(arg0 []byte) (weles.JobID, error) {
	ret := m.ctrl.Call(m, "CreateJob", arg0)
	ret0, _ := ret[0].(weles.JobID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJob indicates an expected call of CreateJob
func (mr *MockJobManagerMockRecorder) CreateJob(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJob", reflect.TypeOf((*MockJobManager)(nil).CreateJob), arg0)
}

// ListJobs mocks base method
func (m *MockJobManager) ListJobs(arg0 []weles.JobID) ([]weles.JobInfo, error) {
	ret := m.ctrl.Call(m, "ListJobs", arg0)
	ret0, _ := ret[0].([]weles.JobInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListJobs indicates an expected call of ListJobs
func (mr *MockJobManagerMockRecorder) ListJobs(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListJobs", reflect.TypeOf((*MockJobManager)(nil).ListJobs), arg0)
}
