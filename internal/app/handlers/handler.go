package handlers

import (
	"github.com/Galish/loyalty-system/internal/app/services"
	"github.com/Galish/loyalty-system/internal/app/services/accrual"
	"github.com/Galish/loyalty-system/internal/app/services/balance"
	"github.com/Galish/loyalty-system/internal/app/services/order"
	"github.com/Galish/loyalty-system/internal/app/services/user"
	"github.com/Galish/loyalty-system/internal/config"
)

type httpHandler struct {
	cfg            *config.Config
	orderService   order.OrderManager
	balanceService balance.BalanceManager
	accrualService accrual.AccrualManager
	userService    user.UserManager
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
