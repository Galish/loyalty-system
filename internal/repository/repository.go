package repository

import (
	"context"
	"errors"

	"github.com/Galish/loyalty-system/internal/model"
)

var (
	ErrOrderExists       = errors.New("order has already been added")
	ErrOrderConflict     = errors.New("order has already been added by another user")
	ErrUserConflict      = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientFunds = errors.New("insufficient funds in the account")
)

type UserRepository interface {
	CreateUser(context.Context, string, string) (*model.User, error)
	GetUserByLogin(context.Context, string) (*model.User, error)
}

type OrderRepository interface {
	CreateOrder(context.Context, *model.Order) error
	UserOrders(context.Context, string) ([]*model.Order, error)
	UpdateOrder(context.Context, *model.Order) error
}

type BalanceRepository interface {
	UserBalance(context.Context, string) (*model.Balance, error)
	Enroll(context.Context, *model.Enrollment) error
	Withdraw(context.Context, *model.Withdrawal) error
	Withdrawals(context.Context, string) ([]*model.Withdrawal, error)
}

type LoyaltyRepository interface {
	OrderRepository
	BalanceRepository
}
