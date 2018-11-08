// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SamsungSLAV/weles/controller (interfaces: Dryader)

// Package mock is a generated GoMock package.
package mock

import (
	weles "github.com/SamsungSLAV/weles"
	notifier "github.com/SamsungSLAV/weles/controller/notifier"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDryader is a mock of Dryader interface
type MockDryader struct {
	ctrl     *gomock.Controller
	recorder *MockDryaderMockRecorder
}

// MockDryaderMockRecorder is the mock recorder for MockDryader
type MockDryaderMockRecorder struct {
	mock *MockDryader
}

// NewMockDryader creates a new mock instance
func NewMockDryader(ctrl *gomock.Controller) *MockDryader {
	mock := &MockDryader{ctrl: ctrl}
	mock.recorder = &MockDryaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDryader) EXPECT() *MockDryaderMockRecorder {
	return m.recorder
}

// CancelJob mocks base method
func (m *MockDryader) CancelJob(arg0 weles.JobID) {
	m.ctrl.Call(m, "CancelJob", arg0)
}

// CancelJob indicates an expected call of CancelJob
func (mr *MockDryaderMockRecorder) CancelJob(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelJob", reflect.TypeOf((*MockDryader)(nil).CancelJob), arg0)
}

// Listen mocks base method
func (m *MockDryader) Listen() <-chan notifier.Notification {
	ret := m.ctrl.Call(m, "Listen")
	ret0, _ := ret[0].(<-chan notifier.Notification)
	return ret0
}

// Listen indicates an expected call of Listen
func (mr *MockDryaderMockRecorder) Listen() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listen", reflect.TypeOf((*MockDryader)(nil).Listen))
}

// SendFail mocks base method
func (m *MockDryader) SendFail(arg0 weles.JobID, arg1 string) {
	m.ctrl.Call(m, "SendFail", arg0, arg1)
}

// SendFail indicates an expected call of SendFail
func (mr *MockDryaderMockRecorder) SendFail(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendFail", reflect.TypeOf((*MockDryader)(nil).SendFail), arg0, arg1)
}

// SendOK mocks base method
func (m *MockDryader) SendOK(arg0 weles.JobID) {
	m.ctrl.Call(m, "SendOK", arg0)
}

// SendOK indicates an expected call of SendOK
func (mr *MockDryaderMockRecorder) SendOK(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendOK", reflect.TypeOf((*MockDryader)(nil).SendOK), arg0)
}

// StartJob mocks base method
func (m *MockDryader) StartJob(arg0 weles.JobID) {
	m.ctrl.Call(m, "StartJob", arg0)
}

// StartJob indicates an expected call of StartJob
func (mr *MockDryaderMockRecorder) StartJob(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartJob", reflect.TypeOf((*MockDryader)(nil).StartJob), arg0)
}