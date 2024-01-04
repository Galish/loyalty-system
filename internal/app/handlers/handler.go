package handlers

import (
	"github.com/Galish/loyalty-system/internal/app/services"
	"github.com/Galish/loyalty-system/internal/app/services/accrual"
	"github.com/Galish/loyalty-system/internal/app/services/auth"
	"github.com/Galish/loyalty-system/internal/app/services/balance"
	"github.com/Galish/loyalty-system/internal/app/services/order"
	"github.com/Galish/loyalty-system/internal/config"
)

type httpHandler struct {
	cfg            *config.Config
	authService    auth.AuthManager
	orderService   order.OrderManager
	balanceService balance.BalanceManager
	accrualService accrual.AccrualManager
}

func NewHandler(
	cfg *config.Config,
	svc *services.Services,
) *httpHandler {
	return &httpHandler{
		cfg:            cfg,
		authService:    svc.Auth,
		orderService:   svc.Order,
		balanceService: svc.Balance,
		accrualService: svc.Accrual,
	}
}
