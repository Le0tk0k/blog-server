// Code generated by MockGen. DO NOT EDIT.
// Source: tag.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "github.com/Le0tk0k/blog-server/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTagRepository is a mock of TagRepository interface
type MockTagRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTagRepositoryMockRecorder
}

// MockTagRepositoryMockRecorder is the mock recorder for MockTagRepository
type MockTagRepositoryMockRecorder struct {
	mock *MockTagRepository
}

// NewMockTagRepository creates a new mock instance
func NewMockTagRepository(ctrl *gomock.Controller) *MockTagRepository {
	mock := &MockTagRepository{ctrl: ctrl}
	mock.recorder = &MockTagRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTagRepository) EXPECT() *MockTagRepositoryMockRecorder {
	return m.recorder
}

// StoreTag mocks base method
func (m *MockTagRepository) StoreTag(tag *model.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreTag", tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreTag indicates an expected call of StoreTag
func (mr *MockTagRepositoryMockRecorder) StoreTag(tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreTag", reflect.TypeOf((*MockTagRepository)(nil).StoreTag), tag)
}
