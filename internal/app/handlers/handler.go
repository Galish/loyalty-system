package handlers

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/app/services/accrual"
	"github.com/Galish/loyalty-system/internal/app/services/auth"
	"github.com/Galish/loyalty-system/internal/app/services/balance"
	"github.com/Galish/loyalty-system/internal/config"
)

type orderService interface {
	AddOrder(context.Context, entity.Order) error
	GetOrders(context.Context, string) ([]*entity.Order, error)
}

type httpHandler struct {
	cfg            *config.Config
	authService    auth.AuthManager
	orderService   orderService
	balanceService balance.BalanceManager
	accrualService accrual.AccrualManager
}

func NewHandler(
	cfg *config.Config,
	auth auth.AuthManager,
	order orderService,
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
