package handlers

import (
	"github.com/Galish/loyalty-system/internal/app/services"
	"github.com/Galish/loyalty-system/internal/config"
)

type httpHandler struct {
	cfg            *config.Config
	orderService   services.OrderManager
	balanceService services.BalanceManager
	accrualService services.AccrualManager
	userService    services.UserManager
}

func NewHandler(
	cfg *config.Config,
	svc *services.Services,
) *httpHandler {
	return &httpHandler{
		accrualService: svc.Accrual,
		balanceService: svc.Balance,
		cfg:            cfg,
		orderService:   svc.Order,
		userService:    svc.User,
	}
}
