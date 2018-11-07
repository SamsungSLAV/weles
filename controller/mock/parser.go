// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SamsungSLAV/weles/controller (interfaces: Parser)

// Package mock is a generated GoMock package.
package mock

import (
	weles "github.com/SamsungSLAV/weles"
	notifier "github.com/SamsungSLAV/weles/controller/notifier"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockParser is a mock of Parser interface
type MockParser struct {
	ctrl     *gomock.Controller
	recorder *MockParserMockRecorder
}

// MockParserMockRecorder is the mock recorder for MockParser
type MockParserMockRecorder struct {
	mock *MockParser
}

// NewMockParser creates a new mock instance
func NewMockParser(ctrl *gomock.Controller) *MockParser {
	mock := &MockParser{ctrl: ctrl}
	mock.recorder = &MockParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockParser) EXPECT() *MockParserMockRecorder {
	return m.recorder
}

// Listen mocks base method
func (m *MockParser) Listen() <-chan notifier.Notification {
	ret := m.ctrl.Call(m, "Listen")
	ret0, _ := ret[0].(<-chan notifier.Notification)
	return ret0
}

// Listen indicates an expected call of Listen
func (mr *MockParserMockRecorder) Listen() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listen", reflect.TypeOf((*MockParser)(nil).Listen))
}

// Parse mocks base method
func (m *MockParser) Parse(arg0 weles.JobID) {
	m.ctrl.Call(m, "Parse", arg0)
}

// Parse indicates an expected call of Parse
func (mr *MockParserMockRecorder) Parse(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockParser)(nil).Parse), arg0)
}

// SendFail mocks base method
func (m *MockParser) SendFail(arg0 weles.JobID, arg1 string) {
	m.ctrl.Call(m, "SendFail", arg0, arg1)
}

// SendFail indicates an expected call of SendFail
func (mr *MockParserMockRecorder) SendFail(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendFail", reflect.TypeOf((*MockParser)(nil).SendFail), arg0, arg1)
}

// SendOK mocks base method
func (m *MockParser) SendOK(arg0 weles.JobID) {
	m.ctrl.Call(m, "SendOK", arg0)
}

// SendOK indicates an expected call of SendOK
func (mr *MockParserMockRecorder) SendOK(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendOK", reflect.TypeOf((*MockParser)(nil).SendOK), arg0)
}