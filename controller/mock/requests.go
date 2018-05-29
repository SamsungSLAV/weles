// Code generated by MockGen. DO NOT EDIT.
// Source: ../../../boruta/boruta.go

// Package mock is a generated GoMock package.
package mock

import (
	rsa "crypto/rsa"
	boruta "git.tizen.org/tools/boruta"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockListFilter is a mock of ListFilter interface
type MockListFilter struct {
	ctrl     *gomock.Controller
	recorder *MockListFilterMockRecorder
}

// MockListFilterMockRecorder is the mock recorder for MockListFilter
type MockListFilterMockRecorder struct {
	mock *MockListFilter
}

// NewMockListFilter creates a new mock instance
func NewMockListFilter(ctrl *gomock.Controller) *MockListFilter {
	mock := &MockListFilter{ctrl: ctrl}
	mock.recorder = &MockListFilterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockListFilter) EXPECT() *MockListFilterMockRecorder {
	return m.recorder
}

// Match mocks base method
func (m *MockListFilter) Match(req *boruta.ReqInfo) bool {
	ret := m.ctrl.Call(m, "Match", req)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Match indicates an expected call of Match
func (mr *MockListFilterMockRecorder) Match(req interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Match", reflect.TypeOf((*MockListFilter)(nil).Match), req)
}

// MockRequests is a mock of Requests interface
type MockRequests struct {
	ctrl     *gomock.Controller
	recorder *MockRequestsMockRecorder
}

// MockRequestsMockRecorder is the mock recorder for MockRequests
type MockRequestsMockRecorder struct {
	mock *MockRequests
}

// NewMockRequests creates a new mock instance
func NewMockRequests(ctrl *gomock.Controller) *MockRequests {
	mock := &MockRequests{ctrl: ctrl}
	mock.recorder = &MockRequestsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequests) EXPECT() *MockRequestsMockRecorder {
	return m.recorder
}

// NewRequest mocks base method
func (m *MockRequests) NewRequest(caps boruta.Capabilities, priority boruta.Priority, owner boruta.UserInfo, validAfter, deadline time.Time) (boruta.ReqID, error) {
	ret := m.ctrl.Call(m, "NewRequest", caps, priority, owner, validAfter, deadline)
	ret0, _ := ret[0].(boruta.ReqID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRequest indicates an expected call of NewRequest
func (mr *MockRequestsMockRecorder) NewRequest(caps, priority, owner, validAfter, deadline interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRequest", reflect.TypeOf((*MockRequests)(nil).NewRequest), caps, priority, owner, validAfter, deadline)
}

// CloseRequest mocks base method
func (m *MockRequests) CloseRequest(reqID boruta.ReqID) error {
	ret := m.ctrl.Call(m, "CloseRequest", reqID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseRequest indicates an expected call of CloseRequest
func (mr *MockRequestsMockRecorder) CloseRequest(reqID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseRequest", reflect.TypeOf((*MockRequests)(nil).CloseRequest), reqID)
}

// UpdateRequest mocks base method
func (m *MockRequests) UpdateRequest(reqInfo *boruta.ReqInfo) error {
	ret := m.ctrl.Call(m, "UpdateRequest", reqInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRequest indicates an expected call of UpdateRequest
func (mr *MockRequestsMockRecorder) UpdateRequest(reqInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRequest", reflect.TypeOf((*MockRequests)(nil).UpdateRequest), reqInfo)
}

// GetRequestInfo mocks base method
func (m *MockRequests) GetRequestInfo(reqID boruta.ReqID) (boruta.ReqInfo, error) {
	ret := m.ctrl.Call(m, "GetRequestInfo", reqID)
	ret0, _ := ret[0].(boruta.ReqInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRequestInfo indicates an expected call of GetRequestInfo
func (mr *MockRequestsMockRecorder) GetRequestInfo(reqID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRequestInfo", reflect.TypeOf((*MockRequests)(nil).GetRequestInfo), reqID)
}

// ListRequests mocks base method
func (m *MockRequests) ListRequests(filter boruta.ListFilter) ([]boruta.ReqInfo, error) {
	ret := m.ctrl.Call(m, "ListRequests", filter)
	ret0, _ := ret[0].([]boruta.ReqInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRequests indicates an expected call of ListRequests
func (mr *MockRequestsMockRecorder) ListRequests(filter interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRequests", reflect.TypeOf((*MockRequests)(nil).ListRequests), filter)
}

// AcquireWorker mocks base method
func (m *MockRequests) AcquireWorker(reqID boruta.ReqID) (boruta.AccessInfo, error) {
	ret := m.ctrl.Call(m, "AcquireWorker", reqID)
	ret0, _ := ret[0].(boruta.AccessInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcquireWorker indicates an expected call of AcquireWorker
func (mr *MockRequestsMockRecorder) AcquireWorker(reqID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcquireWorker", reflect.TypeOf((*MockRequests)(nil).AcquireWorker), reqID)
}

// ProlongAccess mocks base method
func (m *MockRequests) ProlongAccess(reqID boruta.ReqID) error {
	ret := m.ctrl.Call(m, "ProlongAccess", reqID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProlongAccess indicates an expected call of ProlongAccess
func (mr *MockRequestsMockRecorder) ProlongAccess(reqID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProlongAccess", reflect.TypeOf((*MockRequests)(nil).ProlongAccess), reqID)
}

// MockSuperviser is a mock of Superviser interface
type MockSuperviser struct {
	ctrl     *gomock.Controller
	recorder *MockSuperviserMockRecorder
}

// MockSuperviserMockRecorder is the mock recorder for MockSuperviser
type MockSuperviserMockRecorder struct {
	mock *MockSuperviser
}

// NewMockSuperviser creates a new mock instance
func NewMockSuperviser(ctrl *gomock.Controller) *MockSuperviser {
	mock := &MockSuperviser{ctrl: ctrl}
	mock.recorder = &MockSuperviserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSuperviser) EXPECT() *MockSuperviserMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockSuperviser) Register(caps boruta.Capabilities) error {
	ret := m.ctrl.Call(m, "Register", caps)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockSuperviserMockRecorder) Register(caps interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockSuperviser)(nil).Register), caps)
}

// SetFail mocks base method
func (m *MockSuperviser) SetFail(uuid boruta.WorkerUUID, reason string) error {
	ret := m.ctrl.Call(m, "SetFail", uuid, reason)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetFail indicates an expected call of SetFail
func (mr *MockSuperviserMockRecorder) SetFail(uuid, reason interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFail", reflect.TypeOf((*MockSuperviser)(nil).SetFail), uuid, reason)
}

// MockWorkers is a mock of Workers interface
type MockWorkers struct {
	ctrl     *gomock.Controller
	recorder *MockWorkersMockRecorder
}

// MockWorkersMockRecorder is the mock recorder for MockWorkers
type MockWorkersMockRecorder struct {
	mock *MockWorkers
}

// NewMockWorkers creates a new mock instance
func NewMockWorkers(ctrl *gomock.Controller) *MockWorkers {
	mock := &MockWorkers{ctrl: ctrl}
	mock.recorder = &MockWorkersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkers) EXPECT() *MockWorkersMockRecorder {
	return m.recorder
}

// ListWorkers mocks base method
func (m *MockWorkers) ListWorkers(groups boruta.Groups, caps boruta.Capabilities) ([]boruta.WorkerInfo, error) {
	ret := m.ctrl.Call(m, "ListWorkers", groups, caps)
	ret0, _ := ret[0].([]boruta.WorkerInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWorkers indicates an expected call of ListWorkers
func (mr *MockWorkersMockRecorder) ListWorkers(groups, caps interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWorkers", reflect.TypeOf((*MockWorkers)(nil).ListWorkers), groups, caps)
}

// GetWorkerInfo mocks base method
func (m *MockWorkers) GetWorkerInfo(uuid boruta.WorkerUUID) (boruta.WorkerInfo, error) {
	ret := m.ctrl.Call(m, "GetWorkerInfo", uuid)
	ret0, _ := ret[0].(boruta.WorkerInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkerInfo indicates an expected call of GetWorkerInfo
func (mr *MockWorkersMockRecorder) GetWorkerInfo(uuid interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkerInfo", reflect.TypeOf((*MockWorkers)(nil).GetWorkerInfo), uuid)
}

// SetState mocks base method
func (m *MockWorkers) SetState(uuid boruta.WorkerUUID, state boruta.WorkerState) error {
	ret := m.ctrl.Call(m, "SetState", uuid, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetState indicates an expected call of SetState
func (mr *MockWorkersMockRecorder) SetState(uuid, state interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetState", reflect.TypeOf((*MockWorkers)(nil).SetState), uuid, state)
}

// SetGroups mocks base method
func (m *MockWorkers) SetGroups(uuid boruta.WorkerUUID, groups boruta.Groups) error {
	ret := m.ctrl.Call(m, "SetGroups", uuid, groups)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetGroups indicates an expected call of SetGroups
func (mr *MockWorkersMockRecorder) SetGroups(uuid, groups interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGroups", reflect.TypeOf((*MockWorkers)(nil).SetGroups), uuid, groups)
}

// Deregister mocks base method
func (m *MockWorkers) Deregister(uuid boruta.WorkerUUID) error {
	ret := m.ctrl.Call(m, "Deregister", uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deregister indicates an expected call of Deregister
func (mr *MockWorkersMockRecorder) Deregister(uuid interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deregister", reflect.TypeOf((*MockWorkers)(nil).Deregister), uuid)
}

// MockDryad is a mock of Dryad interface
type MockDryad struct {
	ctrl     *gomock.Controller
	recorder *MockDryadMockRecorder
}

// MockDryadMockRecorder is the mock recorder for MockDryad
type MockDryadMockRecorder struct {
	mock *MockDryad
}

// NewMockDryad creates a new mock instance
func NewMockDryad(ctrl *gomock.Controller) *MockDryad {
	mock := &MockDryad{ctrl: ctrl}
	mock.recorder = &MockDryadMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDryad) EXPECT() *MockDryadMockRecorder {
	return m.recorder
}

// PutInMaintenance mocks base method
func (m *MockDryad) PutInMaintenance(msg string) error {
	ret := m.ctrl.Call(m, "PutInMaintenance", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutInMaintenance indicates an expected call of PutInMaintenance
func (mr *MockDryadMockRecorder) PutInMaintenance(msg interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutInMaintenance", reflect.TypeOf((*MockDryad)(nil).PutInMaintenance), msg)
}

// Prepare mocks base method
func (m *MockDryad) Prepare() (*rsa.PrivateKey, error) {
	ret := m.ctrl.Call(m, "Prepare")
	ret0, _ := ret[0].(*rsa.PrivateKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare
func (mr *MockDryadMockRecorder) Prepare() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockDryad)(nil).Prepare))
}

// Healthcheck mocks base method
func (m *MockDryad) Healthcheck() error {
	ret := m.ctrl.Call(m, "Healthcheck")
	ret0, _ := ret[0].(error)
	return ret0
}

// Healthcheck indicates an expected call of Healthcheck
func (mr *MockDryadMockRecorder) Healthcheck() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Healthcheck", reflect.TypeOf((*MockDryad)(nil).Healthcheck))
}
