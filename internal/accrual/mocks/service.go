// Code generated by MockGen. DO NOT EDIT.
// Source: internal/accrual/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/Galish/loyalty-system/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockAccrualManager is a mock of AccrualManager interface.
type MockAccrualManager struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualManagerMockRecorder
}

// MockAccrualManagerMockRecorder is the mock recorder for MockAccrualManager.
type MockAccrualManagerMockRecorder struct {
	mock *MockAccrualManager
}

// NewMockAccrualManager creates a new mock instance.
func NewMockAccrualManager(ctrl *gomock.Controller) *MockAccrualManager {
	mock := &MockAccrualManager{ctrl: ctrl}
	mock.recorder = &MockAccrualManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualManager) EXPECT() *MockAccrualManagerMockRecorder {
	return m.recorder
}

// GetAccrual mocks base method.
func (m *MockAccrualManager) GetAccrual(order *model.Order) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetAccrual", order)
}

// GetAccrual indicates an expected call of GetAccrual.
func (mr *MockAccrualManagerMockRecorder) GetAccrual(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccrual", reflect.TypeOf((*MockAccrualManager)(nil).GetAccrual), order)
}