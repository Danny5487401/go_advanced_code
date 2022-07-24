// Code generated by MockGen. DO NOT EDIT.
// Source: dao.go

// Package dao is a generated GoMock package.
package dao

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSearch is a mock of Search interface.
type MockSearch struct {
	ctrl     *gomock.Controller
	recorder *MockSearchMockRecorder
}

// MockSearchMockRecorder is the mock recorder for MockSearch.
type MockSearchMockRecorder struct {
	mock *MockSearch
}

// NewMockSearch creates a new mock instance.
func NewMockSearch(ctrl *gomock.Controller) *MockSearch {
	mock := &MockSearch{ctrl: ctrl}
	mock.recorder = &MockSearchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearch) EXPECT() *MockSearchMockRecorder {
	return m.recorder
}

// GetNameByID mocks base method.
func (m *MockSearch) GetNameByID(id int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNameByID", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNameByID indicates an expected call of GetNameByID.
func (mr *MockSearchMockRecorder) GetNameByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNameByID", reflect.TypeOf((*MockSearch)(nil).GetNameByID), id)
}
