package api

import (
	"github.com/Galish/loyalty-system/internal/accrual"
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/balance"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/order"
)

type httpHandler struct {
	cfg            *config.Config
	authService    auth.AuthManager
	orderService   order.OrderManager
	balanceService balance.BalanceManager
	accrualService accrual.AccrualManager
}

func newHandler(
	cfg *config.Config,
	auth auth.AuthManager,
	order order.OrderManager,
	balance balance.BalanceManager,
	accrual accrual.AccrualManager,
) *httpHandler {
	return &httpHandler{
		cfg:            cfg,
		authService:    auth,
		orderService:   order,
		balanceService: balance,
		accrualService: accrual,
	}
}
