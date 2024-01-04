package repository

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

type UserRepository interface {
	CreateUser(context.Context, string, string) (*entity.User, error)
	GetUserByLogin(context.Context, string) (*entity.User, error)
}

type OrderRepository interface {
	CreateOrder(context.Context, *entity.Order) error
	UserOrders(context.Context, string) ([]*entity.Order, error)
	UpdateOrder(context.Context, *entity.Order) error
}

type BalanceRepository interface {
	UserBalance(context.Context, string) (*entity.Balance, error)
	Enroll(context.Context, *entity.Enrollment) error
	Withdraw(context.Context, *entity.Withdrawal) error
	Withdrawals(context.Context, string) ([]*entity.Withdrawal, error)
}

type Repository interface {
	UserRepository
	OrderRepository
	BalanceRepository
}
