package handlers

import (
	"github.com/Galish/loyalty-system/internal/app/services"
	"github.com/Galish/loyalty-system/internal/config"
)

type httpHandler struct {
	cfg     *config.Config
	order   services.OrderManager
	balance services.BalanceManager
	accrual services.AccrualManager
	user    services.UserManager
}

func NewHandler(
	cfg *config.Config,
	svc *services.Services,
) *httpHandler {
	return &httpHandler{
		cfg:     cfg,
		accrual: svc.Accrual,
		balance: svc.Balance,
		order:   svc.Order,
		user:    svc.User,
	}
}
