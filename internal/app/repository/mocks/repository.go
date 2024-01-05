// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/repository/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/Galish/loyalty-system/internal/app/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(arg0 context.Context, arg1, arg2 string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), arg0, arg1, arg2)
}

// GetUserByLogin mocks base method.
func (m *MockUserRepository) GetUserByLogin(arg0 context.Context, arg1 string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockUserRepositoryMockRecorder) GetUserByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockUserRepository)(nil).GetUserByLogin), arg0, arg1)
}

// MockOrderRepository is a mock of OrderRepository interface.
type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
}

// MockOrderRepositoryMockRecorder is the mock recorder for MockOrderRepository.
type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

// NewMockOrderRepository creates a new mock instance.
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrderRepository) CreateOrder(arg0 context.Context, arg1 *entity.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderRepositoryMockRecorder) CreateOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderRepository)(nil).CreateOrder), arg0, arg1)
}

// UpdateOrder mocks base method.
func (m *MockOrderRepository) UpdateOrder(arg0 context.Context, arg1 *entity.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockOrderRepositoryMockRecorder) UpdateOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockOrderRepository)(nil).UpdateOrder), arg0, arg1)
}

// UserOrders mocks base method.
func (m *MockOrderRepository) UserOrders(arg0 context.Context, arg1 string) ([]*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserOrders", arg0, arg1)
	ret0, _ := ret[0].([]*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserOrders indicates an expected call of UserOrders.
func (mr *MockOrderRepositoryMockRecorder) UserOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserOrders", reflect.TypeOf((*MockOrderRepository)(nil).UserOrders), arg0, arg1)
}

// MockBalanceRepository is a mock of BalanceRepository interface.
type MockBalanceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceRepositoryMockRecorder
}

// MockBalanceRepositoryMockRecorder is the mock recorder for MockBalanceRepository.
type MockBalanceRepositoryMockRecorder struct {
	mock *MockBalanceRepository
}

// NewMockBalanceRepository creates a new mock instance.
func NewMockBalanceRepository(ctrl *gomock.Controller) *MockBalanceRepository {
	mock := &MockBalanceRepository{ctrl: ctrl}
	mock.recorder = &MockBalanceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalanceRepository) EXPECT() *MockBalanceRepositoryMockRecorder {
	return m.recorder
}

// Enroll mocks base method.
func (m *MockBalanceRepository) Enroll(arg0 context.Context, arg1 *entity.Enrollment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enroll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enroll indicates an expected call of Enroll.
func (mr *MockBalanceRepositoryMockRecorder) Enroll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enroll", reflect.TypeOf((*MockBalanceRepository)(nil).Enroll), arg0, arg1)
}

// UserBalance mocks base method.
func (m *MockBalanceRepository) UserBalance(arg0 context.Context, arg1 string) (*entity.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserBalance", arg0, arg1)
	ret0, _ := ret[0].(*entity.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserBalance indicates an expected call of UserBalance.
func (mr *MockBalanceRepositoryMockRecorder) UserBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserBalance", reflect.TypeOf((*MockBalanceRepository)(nil).UserBalance), arg0, arg1)
}

// Withdraw mocks base method.
func (m *MockBalanceRepository) Withdraw(arg0 context.Context, arg1 *entity.Withdrawal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockBalanceRepositoryMockRecorder) Withdraw(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockBalanceRepository)(nil).Withdraw), arg0, arg1)
}

// Withdrawals mocks base method.
func (m *MockBalanceRepository) Withdrawals(arg0 context.Context, arg1 string) ([]*entity.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdrawals", arg0, arg1)
	ret0, _ := ret[0].([]*entity.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Withdrawals indicates an expected call of Withdrawals.
func (mr *MockBalanceRepositoryMockRecorder) Withdrawals(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdrawals", reflect.TypeOf((*MockBalanceRepository)(nil).Withdrawals), arg0, arg1)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockRepository) CreateOrder(arg0 context.Context, arg1 *entity.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockRepositoryMockRecorder) CreateOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockRepository)(nil).CreateOrder), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(arg0 context.Context, arg1, arg2 string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), arg0, arg1, arg2)
}

// Enroll mocks base method.
func (m *MockRepository) Enroll(arg0 context.Context, arg1 *entity.Enrollment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enroll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enroll indicates an expected call of Enroll.
func (mr *MockRepositoryMockRecorder) Enroll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enroll", reflect.TypeOf((*MockRepository)(nil).Enroll), arg0, arg1)
}

// GetUserByLogin mocks base method.
func (m *MockRepository) GetUserByLogin(arg0 context.Context, arg1 string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockRepositoryMockRecorder) GetUserByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockRepository)(nil).GetUserByLogin), arg0, arg1)
}

// UpdateOrder mocks base method.
func (m *MockRepository) UpdateOrder(arg0 context.Context, arg1 *entity.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockRepositoryMockRecorder) UpdateOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockRepository)(nil).UpdateOrder), arg0, arg1)
}

// UserBalance mocks base method.
func (m *MockRepository) UserBalance(arg0 context.Context, arg1 string) (*entity.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserBalance", arg0, arg1)
	ret0, _ := ret[0].(*entity.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserBalance indicates an expected call of UserBalance.
func (mr *MockRepositoryMockRecorder) UserBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserBalance", reflect.TypeOf((*MockRepository)(nil).UserBalance), arg0, arg1)
}

// UserOrders mocks base method.
func (m *MockRepository) UserOrders(arg0 context.Context, arg1 string) ([]*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserOrders", arg0, arg1)
	ret0, _ := ret[0].([]*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserOrders indicates an expected call of UserOrders.
func (mr *MockRepositoryMockRecorder) UserOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserOrders", reflect.TypeOf((*MockRepository)(nil).UserOrders), arg0, arg1)
}

// Withdraw mocks base method.
func (m *MockRepository) Withdraw(arg0 context.Context, arg1 *entity.Withdrawal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockRepositoryMockRecorder) Withdraw(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockRepository)(nil).Withdraw), arg0, arg1)
}

// Withdrawals mocks base method.
func (m *MockRepository) Withdrawals(arg0 context.Context, arg1 string) ([]*entity.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdrawals", arg0, arg1)
	ret0, _ := ret[0].([]*entity.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Withdrawals indicates an expected call of Withdrawals.
func (mr *MockRepositoryMockRecorder) Withdrawals(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdrawals", reflect.TypeOf((*MockRepository)(nil).Withdrawals), arg0, arg1)
}
