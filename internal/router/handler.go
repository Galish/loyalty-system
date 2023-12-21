package router

import (
	"net/http"

	"github.com/Galish/loyalty-system/internal/accrual"
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/balance"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/order"
)

type httpHandler struct {
	cfg            *config.Config
	authService    *auth.AuthService
	orderService   *order.OrderService
	balanceService *balance.BalanceService
	accrualService accrual.AccrualManager
}

func newHandler(
	cfg *config.Config,
	auth *auth.AuthService,
	order *order.OrderService,
	balance *balance.BalanceService,
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

func (h *httpHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
