// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/artkescha/checker/online_checker/pkg/zipper (interfaces: Zipper)

// Package zipper is a generated GoMock package.
package zipper

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockZipper is a mock of Zipper interface.
type MockZipper struct {
	ctrl     *gomock.Controller
	recorder *MockZipperMockRecorder
}

// MockZipperMockRecorder is the mock recorder for MockZipper.
type MockZipperMockRecorder struct {
	mock *MockZipper
}

// NewMockZipper creates a new mock instance.
func NewMockZipper(ctrl *gomock.Controller) *MockZipper {
	mock := &MockZipper{ctrl: ctrl}
	mock.recorder = &MockZipperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockZipper) EXPECT() *MockZipperMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockZipper) Add(arg0 []string, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockZipperMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockZipper)(nil).Add), arg0, arg1)
}

// Get mocks base method.
func (m *MockZipper) Get() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockZipperMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockZipper)(nil).Get))
}
