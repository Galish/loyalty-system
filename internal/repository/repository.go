package repository

import (
	"context"

	"github.com/Galish/loyalty-system/internal/model"
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

// type LoyaltyRepository interface {
// 	OrderRepository
// 	BalanceRepository
// }
