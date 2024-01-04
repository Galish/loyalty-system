package services

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
	"github.com/Galish/loyalty-system/internal/app/services/accrual"
	"github.com/Galish/loyalty-system/internal/app/services/balance"
	"github.com/Galish/loyalty-system/internal/app/services/order"
	"github.com/Galish/loyalty-system/internal/app/services/user"
	"github.com/Galish/loyalty-system/internal/app/webapi"
	"github.com/Galish/loyalty-system/internal/config"
)

type Services struct {
	Accrual accrual.AccrualManager
	Balance balance.BalanceManager
	Order   order.OrderManager
	User    user.UserManager
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
