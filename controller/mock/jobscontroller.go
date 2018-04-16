// Code generated by MockGen. DO NOT EDIT.
// Source: git.tizen.org/tools/weles/controller (interfaces: JobsController)

// Package mock is a generated GoMock package.
package mock

import (
	weles "git.tizen.org/tools/weles"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockJobsController is a mock of JobsController interface
type MockJobsController struct {
	ctrl     *gomock.Controller
	recorder *MockJobsControllerMockRecorder
}

// MockJobsControllerMockRecorder is the mock recorder for MockJobsController
type MockJobsControllerMockRecorder struct {
	mock *MockJobsController
}

// NewMockJobsController creates a new mock instance
func NewMockJobsController(ctrl *gomock.Controller) *MockJobsController {
	mock := &MockJobsController{ctrl: ctrl}
	mock.recorder = &MockJobsControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJobsController) EXPECT() *MockJobsControllerMockRecorder {
	return m.recorder
}

// GetConfig mocks base method
func (m *MockJobsController) GetConfig(arg0 weles.JobID) (weles.Config, error) {
	ret := m.ctrl.Call(m, "GetConfig", arg0)
	ret0, _ := ret[0].(weles.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfig indicates an expected call of GetConfig
func (mr *MockJobsControllerMockRecorder) GetConfig(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockJobsController)(nil).GetConfig), arg0)
}

// GetDryad mocks base method
func (m *MockJobsController) GetDryad(arg0 weles.JobID) (weles.Dryad, error) {
	ret := m.ctrl.Call(m, "GetDryad", arg0)
	ret0, _ := ret[0].(weles.Dryad)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDryad indicates an expected call of GetDryad
func (mr *MockJobsControllerMockRecorder) GetDryad(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDryad", reflect.TypeOf((*MockJobsController)(nil).GetDryad), arg0)
}

// GetYaml mocks base method
func (m *MockJobsController) GetYaml(arg0 weles.JobID) ([]byte, error) {
	ret := m.ctrl.Call(m, "GetYaml", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetYaml indicates an expected call of GetYaml
func (mr *MockJobsControllerMockRecorder) GetYaml(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetYaml", reflect.TypeOf((*MockJobsController)(nil).GetYaml), arg0)
}

// List mocks base method
func (m *MockJobsController) List(arg0 []weles.JobID) ([]weles.JobInfo, error) {
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]weles.JobInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockJobsControllerMockRecorder) List(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockJobsController)(nil).List), arg0)
}

// NewJob mocks base method
func (m *MockJobsController) NewJob(arg0 []byte) (weles.JobID, error) {
	ret := m.ctrl.Call(m, "NewJob", arg0)
	ret0, _ := ret[0].(weles.JobID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewJob indicates an expected call of NewJob
func (mr *MockJobsControllerMockRecorder) NewJob(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewJob", reflect.TypeOf((*MockJobsController)(nil).NewJob), arg0)
}

// SetConfig mocks base method
func (m *MockJobsController) SetConfig(arg0 weles.JobID, arg1 weles.Config) error {
	ret := m.ctrl.Call(m, "SetConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConfig indicates an expected call of SetConfig
func (mr *MockJobsControllerMockRecorder) SetConfig(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfig", reflect.TypeOf((*MockJobsController)(nil).SetConfig), arg0, arg1)
}

// SetDryad mocks base method
func (m *MockJobsController) SetDryad(arg0 weles.JobID, arg1 weles.Dryad) error {
	ret := m.ctrl.Call(m, "SetDryad", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDryad indicates an expected call of SetDryad
func (mr *MockJobsControllerMockRecorder) SetDryad(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDryad", reflect.TypeOf((*MockJobsController)(nil).SetDryad), arg0, arg1)
}

// SetStatusAndInfo mocks base method
func (m *MockJobsController) SetStatusAndInfo(arg0 weles.JobID, arg1 weles.JobStatus, arg2 string) error {
	ret := m.ctrl.Call(m, "SetStatusAndInfo", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStatusAndInfo indicates an expected call of SetStatusAndInfo
func (mr *MockJobsControllerMockRecorder) SetStatusAndInfo(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatusAndInfo", reflect.TypeOf((*MockJobsController)(nil).SetStatusAndInfo), arg0, arg1, arg2)
}
