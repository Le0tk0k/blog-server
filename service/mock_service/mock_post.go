// Code generated by MockGen. DO NOT EDIT.
// Source: post.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	model "github.com/Le0tk0k/blog-server/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPostService is a mock of PostService interface
type MockPostService struct {
	ctrl     *gomock.Controller
	recorder *MockPostServiceMockRecorder
}

// MockPostServiceMockRecorder is the mock recorder for MockPostService
type MockPostServiceMockRecorder struct {
	mock *MockPostService
}

// NewMockPostService creates a new mock instance
func NewMockPostService(ctrl *gomock.Controller) *MockPostService {
	mock := &MockPostService{ctrl: ctrl}
	mock.recorder = &MockPostServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPostService) EXPECT() *MockPostServiceMockRecorder {
	return m.recorder
}

// CreatePost mocks base method
func (m *MockPostService) CreatePost() (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost")
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost
func (mr *MockPostServiceMockRecorder) CreatePost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockPostService)(nil).CreatePost))
}

// GetPost mocks base method
func (m *MockPostService) GetPost(slug string) (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPost", slug)
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPost indicates an expected call of GetPost
func (mr *MockPostServiceMockRecorder) GetPost(slug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPost", reflect.TypeOf((*MockPostService)(nil).GetPost), slug)
}

// GetPosts mocks base method
func (m *MockPostService) GetPosts(conditions []string) ([]*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPosts", conditions)
	ret0, _ := ret[0].([]*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPosts indicates an expected call of GetPosts
func (mr *MockPostServiceMockRecorder) GetPosts(conditions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPosts", reflect.TypeOf((*MockPostService)(nil).GetPosts), conditions)
}

// UpdatePost mocks base method
func (m *MockPostService) UpdatePost(post *model.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePost", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePost indicates an expected call of UpdatePost
func (mr *MockPostServiceMockRecorder) UpdatePost(post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePost", reflect.TypeOf((*MockPostService)(nil).UpdatePost), post)
}

// DeletePost mocks base method
func (m *MockPostService) DeletePost(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost
func (mr *MockPostServiceMockRecorder) DeletePost(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockPostService)(nil).DeletePost), id)
}
