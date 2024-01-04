package services

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/app/repository"
	"github.com/Galish/loyalty-system/internal/app/services/accrual"
	"github.com/Galish/loyalty-system/internal/app/services/balance"
	"github.com/Galish/loyalty-system/internal/app/services/order"
	"github.com/Galish/loyalty-system/internal/app/services/user"
	"github.com/Galish/loyalty-system/internal/app/webapi"
	"github.com/Galish/loyalty-system/internal/config"
)

type AccrualManager interface {
	GetAccrual(context.Context, *entity.Order)
	Close()
}

type BalanceManager interface {
	GetBalance(context.Context, string) (*entity.Balance, error)
	Withdraw(context.Context, *entity.Withdrawal) error
	Withdrawals(context.Context, string) ([]*entity.Withdrawal, error)
}

type OrderManager interface {
	AddOrder(context.Context, entity.Order) error
	GetOrders(context.Context, string) ([]*entity.Order, error)
}

type UserManager interface {
	Register(context.Context, string, string) (string, error)
	Authenticate(context.Context, string, string) (string, error)
}

type Services struct {
	Accrual AccrualManager
	Balance BalanceManager
	Order   OrderManager
	User    UserManager
}

func New(
	cfg *config.Config,
	store repository.Repository,
	webAPI webapi.AccrualGetter,
) *Services {
	return &Services{
		Accrual: accrual.New(webAPI, store, store, cfg),
		Balance: balance.New(store),
		Order:   order.New(store),
		User:    user.New(store, cfg.SecretKey),
	}
}

func (s *Services) Close() {
	s.Accrual.Close()
}
