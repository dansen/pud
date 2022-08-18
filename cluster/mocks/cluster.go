// Code generated by MockGen. DO NOT EDIT.
// Source: cluster/cluster.go

// Package mock_cluster is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cluster "github.com/dansen/pud/cluster"
	message "github.com/dansen/pud/conn/message"
	protos "github.com/dansen/pud/protos"
	route "github.com/dansen/pud/route"
	session "github.com/dansen/pud/session"
)

// MockRPCServer is a mock of RPCServer interface.
type MockRPCServer struct {
	ctrl     *gomock.Controller
	recorder *MockRPCServerMockRecorder
}

// MockRPCServerMockRecorder is the mock recorder for MockRPCServer.
type MockRPCServerMockRecorder struct {
	mock *MockRPCServer
}

// NewMockRPCServer creates a new mock instance.
func NewMockRPCServer(ctrl *gomock.Controller) *MockRPCServer {
	mock := &MockRPCServer{ctrl: ctrl}
	mock.recorder = &MockRPCServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRPCServer) EXPECT() *MockRPCServerMockRecorder {
	return m.recorder
}

// AfterInit mocks base method.
func (m *MockRPCServer) AfterInit() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AfterInit")
}

// AfterInit indicates an expected call of AfterInit.
func (mr *MockRPCServerMockRecorder) AfterInit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterInit", reflect.TypeOf((*MockRPCServer)(nil).AfterInit))
}

// BeforeShutdown mocks base method.
func (m *MockRPCServer) BeforeShutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BeforeShutdown")
}

// BeforeShutdown indicates an expected call of BeforeShutdown.
func (mr *MockRPCServerMockRecorder) BeforeShutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeShutdown", reflect.TypeOf((*MockRPCServer)(nil).BeforeShutdown))
}

// Init mocks base method.
func (m *MockRPCServer) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockRPCServerMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockRPCServer)(nil).Init))
}

// SetPitayaServer mocks base method.
func (m *MockRPCServer) SetPitayaServer(arg0 protos.PitayaServer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPitayaServer", arg0)
}

// SetPitayaServer indicates an expected call of SetPitayaServer.
func (mr *MockRPCServerMockRecorder) SetPitayaServer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPitayaServer", reflect.TypeOf((*MockRPCServer)(nil).SetPitayaServer), arg0)
}

// Shutdown mocks base method.
func (m *MockRPCServer) Shutdown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockRPCServerMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockRPCServer)(nil).Shutdown))
}

// MockRPCClient is a mock of RPCClient interface.
type MockRPCClient struct {
	ctrl     *gomock.Controller
	recorder *MockRPCClientMockRecorder
}

// MockRPCClientMockRecorder is the mock recorder for MockRPCClient.
type MockRPCClientMockRecorder struct {
	mock *MockRPCClient
}

// NewMockRPCClient creates a new mock instance.
func NewMockRPCClient(ctrl *gomock.Controller) *MockRPCClient {
	mock := &MockRPCClient{ctrl: ctrl}
	mock.recorder = &MockRPCClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRPCClient) EXPECT() *MockRPCClientMockRecorder {
	return m.recorder
}

// AfterInit mocks base method.
func (m *MockRPCClient) AfterInit() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AfterInit")
}

// AfterInit indicates an expected call of AfterInit.
func (mr *MockRPCClientMockRecorder) AfterInit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterInit", reflect.TypeOf((*MockRPCClient)(nil).AfterInit))
}

// BeforeShutdown mocks base method.
func (m *MockRPCClient) BeforeShutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BeforeShutdown")
}

// BeforeShutdown indicates an expected call of BeforeShutdown.
func (mr *MockRPCClientMockRecorder) BeforeShutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeShutdown", reflect.TypeOf((*MockRPCClient)(nil).BeforeShutdown))
}

// BroadcastSessionBind mocks base method.
func (m *MockRPCClient) BroadcastSessionBind(uid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BroadcastSessionBind", uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// BroadcastSessionBind indicates an expected call of BroadcastSessionBind.
func (mr *MockRPCClientMockRecorder) BroadcastSessionBind(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastSessionBind", reflect.TypeOf((*MockRPCClient)(nil).BroadcastSessionBind), uid)
}

// Call mocks base method.
func (m *MockRPCClient) Call(ctx context.Context, rpcType protos.RPCType, route *route.Route, session session.Session, msg *message.Message, server *cluster.Server) (*protos.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, rpcType, route, session, msg, server)
	ret0, _ := ret[0].(*protos.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call.
func (mr *MockRPCClientMockRecorder) Call(ctx, rpcType, route, session, msg, server interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockRPCClient)(nil).Call), ctx, rpcType, route, session, msg, server)
}

// Init mocks base method.
func (m *MockRPCClient) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockRPCClientMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockRPCClient)(nil).Init))
}

// Send mocks base method.
func (m *MockRPCClient) Send(route string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", route, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockRPCClientMockRecorder) Send(route, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockRPCClient)(nil).Send), route, data)
}

// SendKick mocks base method.
func (m *MockRPCClient) SendKick(userID, serverType string, kick *protos.KickMsg) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendKick", userID, serverType, kick)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendKick indicates an expected call of SendKick.
func (mr *MockRPCClientMockRecorder) SendKick(userID, serverType, kick interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendKick", reflect.TypeOf((*MockRPCClient)(nil).SendKick), userID, serverType, kick)
}

// SendPush mocks base method.
func (m *MockRPCClient) SendPush(userID string, frontendSv *cluster.Server, push *protos.Push) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendPush", userID, frontendSv, push)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendPush indicates an expected call of SendPush.
func (mr *MockRPCClientMockRecorder) SendPush(userID, frontendSv, push interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendPush", reflect.TypeOf((*MockRPCClient)(nil).SendPush), userID, frontendSv, push)
}

// Shutdown mocks base method.
func (m *MockRPCClient) Shutdown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockRPCClientMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockRPCClient)(nil).Shutdown))
}

// MockSDListener is a mock of SDListener interface.
type MockSDListener struct {
	ctrl     *gomock.Controller
	recorder *MockSDListenerMockRecorder
}

// MockSDListenerMockRecorder is the mock recorder for MockSDListener.
type MockSDListenerMockRecorder struct {
	mock *MockSDListener
}

// NewMockSDListener creates a new mock instance.
func NewMockSDListener(ctrl *gomock.Controller) *MockSDListener {
	mock := &MockSDListener{ctrl: ctrl}
	mock.recorder = &MockSDListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSDListener) EXPECT() *MockSDListenerMockRecorder {
	return m.recorder
}

// AddServer mocks base method.
func (m *MockSDListener) AddServer(arg0 *cluster.Server) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddServer", arg0)
}

// AddServer indicates an expected call of AddServer.
func (mr *MockSDListenerMockRecorder) AddServer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddServer", reflect.TypeOf((*MockSDListener)(nil).AddServer), arg0)
}

// RemoveServer mocks base method.
func (m *MockSDListener) RemoveServer(arg0 *cluster.Server) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveServer", arg0)
}

// RemoveServer indicates an expected call of RemoveServer.
func (mr *MockSDListenerMockRecorder) RemoveServer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServer", reflect.TypeOf((*MockSDListener)(nil).RemoveServer), arg0)
}

// MockRemoteBindingListener is a mock of RemoteBindingListener interface.
type MockRemoteBindingListener struct {
	ctrl     *gomock.Controller
	recorder *MockRemoteBindingListenerMockRecorder
}

// MockRemoteBindingListenerMockRecorder is the mock recorder for MockRemoteBindingListener.
type MockRemoteBindingListenerMockRecorder struct {
	mock *MockRemoteBindingListener
}

// NewMockRemoteBindingListener creates a new mock instance.
func NewMockRemoteBindingListener(ctrl *gomock.Controller) *MockRemoteBindingListener {
	mock := &MockRemoteBindingListener{ctrl: ctrl}
	mock.recorder = &MockRemoteBindingListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRemoteBindingListener) EXPECT() *MockRemoteBindingListenerMockRecorder {
	return m.recorder
}

// OnUserBind mocks base method.
func (m *MockRemoteBindingListener) OnUserBind(uid, fid string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnUserBind", uid, fid)
}

// OnUserBind indicates an expected call of OnUserBind.
func (mr *MockRemoteBindingListenerMockRecorder) OnUserBind(uid, fid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnUserBind", reflect.TypeOf((*MockRemoteBindingListener)(nil).OnUserBind), uid, fid)
}

// MockInfoRetriever is a mock of InfoRetriever interface.
type MockInfoRetriever struct {
	ctrl     *gomock.Controller
	recorder *MockInfoRetrieverMockRecorder
}

// MockInfoRetrieverMockRecorder is the mock recorder for MockInfoRetriever.
type MockInfoRetrieverMockRecorder struct {
	mock *MockInfoRetriever
}

// NewMockInfoRetriever creates a new mock instance.
func NewMockInfoRetriever(ctrl *gomock.Controller) *MockInfoRetriever {
	mock := &MockInfoRetriever{ctrl: ctrl}
	mock.recorder = &MockInfoRetrieverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInfoRetriever) EXPECT() *MockInfoRetrieverMockRecorder {
	return m.recorder
}

// Region mocks base method.
func (m *MockInfoRetriever) Region() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Region")
	ret0, _ := ret[0].(string)
	return ret0
}

// Region indicates an expected call of Region.
func (mr *MockInfoRetrieverMockRecorder) Region() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Region", reflect.TypeOf((*MockInfoRetriever)(nil).Region))
}
