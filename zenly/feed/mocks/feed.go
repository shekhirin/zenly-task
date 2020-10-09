// Code generated by MockGen. DO NOT EDIT.
// Source: feed.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	pb "github.com/shekhirin/zenly-task/zenly/pb"
	reflect "reflect"
)

// MockFeed is a mock of Feed interface
type MockFeed struct {
	ctrl     *gomock.Controller
	recorder *MockFeedMockRecorder
}

// MockFeedMockRecorder is the mock recorder for MockFeed
type MockFeedMockRecorder struct {
	mock *MockFeed
}

// NewMockFeed creates a new mock instance
func NewMockFeed(ctrl *gomock.Controller) *MockFeed {
	mock := &MockFeed{ctrl: ctrl}
	mock.recorder = &MockFeedMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFeed) EXPECT() *MockFeedMockRecorder {
	return m.recorder
}

// Publish mocks base method
func (m *MockFeed) Publish(message *pb.FeedMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockFeedMockRecorder) Publish(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockFeed)(nil).Publish), message)
}