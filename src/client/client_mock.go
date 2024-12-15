// Code generated by MockGen. DO NOT EDIT.
// Source: src/client/client.go

// Package client is a generated GoMock package.
package client

import (
	context "context"
	reflect "reflect"

	types "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	gomock "github.com/golang/mock/gomock"
)

// MockECSClient is a mock of ECSClient interface.
type MockECSClient struct {
	ctrl     *gomock.Controller
	recorder *MockECSClientMockRecorder
}

// MockECSClientMockRecorder is the mock recorder for MockECSClient.
type MockECSClientMockRecorder struct {
	mock *MockECSClient
}

// NewMockECSClient creates a new mock instance.
func NewMockECSClient(ctrl *gomock.Controller) *MockECSClient {
	mock := &MockECSClient{ctrl: ctrl}
	mock.recorder = &MockECSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockECSClient) EXPECT() *MockECSClientMockRecorder {
	return m.recorder
}

// DescribeServices mocks base method.
func (m *MockECSClient) DescribeServices(ctx context.Context, cluster string) ([]types.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeServices", ctx, cluster)
	ret0, _ := ret[0].([]types.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeServices indicates an expected call of DescribeServices.
func (mr *MockECSClientMockRecorder) DescribeServices(ctx, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeServices", reflect.TypeOf((*MockECSClient)(nil).DescribeServices), ctx, cluster)
}

// ListClusters mocks base method.
func (m *MockECSClient) ListClusters(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClusters", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClusters indicates an expected call of ListClusters.
func (mr *MockECSClientMockRecorder) ListClusters(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusters", reflect.TypeOf((*MockECSClient)(nil).ListClusters), ctx)
}

// ScaleinService mocks base method.
func (m *MockECSClient) ScaleinService(ctx context.Context, cluster, service string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScaleinService", ctx, cluster, service)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScaleinService indicates an expected call of ScaleinService.
func (mr *MockECSClientMockRecorder) ScaleinService(ctx, cluster, service interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScaleinService", reflect.TypeOf((*MockECSClient)(nil).ScaleinService), ctx, cluster, service)
}
