package services

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
	"github.com/Galish/loyalty-system/internal/app/services/accrual"
	"github.com/Galish/loyalty-system/internal/app/services/auth"
	"github.com/Galish/loyalty-system/internal/app/services/balance"
	"github.com/Galish/loyalty-system/internal/app/services/order"
	"github.com/Galish/loyalty-system/internal/app/webapi"
	"github.com/Galish/loyalty-system/internal/config"
)

type Services struct {
	Accrual accrual.AccrualManager
	Auth    auth.AuthManager
	Balance balance.BalanceManager
	Order   order.OrderManager
}

func New(
	cfg *config.Config,
	store repository.Repository,
	webAPI webapi.AccrualGetter,
) *Services {
	return &Services{
		Accrual: accrual.New(webAPI, store, store, cfg),
		Auth:    auth.New(store, cfg.SecretKey),
		Balance: balance.New(store),
		Order:   order.New(store),
	}
}

func (s *Services) Close() {
	s.Accrual.Close()
}
