// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/artkescha/checker/online_checker/pkg/unzipper (interfaces: UnZipper)

// Package unzipper is a generated GoMock package.
package unzipper

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUnZipper is a mock of UnZipper interface.
type MockUnZipper struct {
	ctrl     *gomock.Controller
	recorder *MockUnZipperMockRecorder
}

// MockUnZipperMockRecorder is the mock recorder for MockUnZipper.
type MockUnZipperMockRecorder struct {
	mock *MockUnZipper
}

// NewMockUnZipper creates a new mock instance.
func NewMockUnZipper(ctrl *gomock.Controller) *MockUnZipper {
	mock := &MockUnZipper{ctrl: ctrl}
	mock.recorder = &MockUnZipperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnZipper) EXPECT() *MockUnZipperMockRecorder {
	return m.recorder
}

// Unzip mocks base method.
func (m *MockUnZipper) Unzip(arg0, arg1 string, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unzip", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unzip indicates an expected call of Unzip.
func (mr *MockUnZipperMockRecorder) Unzip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unzip", reflect.TypeOf((*MockUnZipper)(nil).Unzip), arg0, arg1, arg2)
}