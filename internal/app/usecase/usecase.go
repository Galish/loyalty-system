package usecase

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

type AccrualUseCase interface {
	GetAccrual(context.Context, *entity.Order)
}

type BalanceUseCase interface {
	GetBalance(context.Context, string) (*entity.Balance, error)
	Withdraw(context.Context, *entity.Withdrawal) error
	Withdrawals(context.Context, string) ([]*entity.Withdrawal, error)
}

type OrderUseCase interface {
	AddOrder(context.Context, entity.Order) error
	GetOrders(context.Context, string) ([]*entity.Order, error)
}

type UserUseCase interface {
	Register(context.Context, string, string) (string, error)
	Authenticate(context.Context, string, string) (string, error)
}
