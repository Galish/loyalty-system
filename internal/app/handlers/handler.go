package handlers

import (
	"github.com/Galish/loyalty-system/internal/app/usecase"
	"github.com/Galish/loyalty-system/internal/config"
)

type httpHandler struct {
	cfg *config.Config
	uc  *useCase
}

type useCase struct {
	accrual usecase.AccrualUseCase
	balance usecase.BalanceUseCase
	order   usecase.OrderUseCase
	user    usecase.UserUseCase
}

func NewHandler(
	cfg *config.Config,
	accrualUseCase usecase.AccrualUseCase,
	balanceUseCase usecase.BalanceUseCase,
	orderUseCase usecase.OrderUseCase,
	userUseCase usecase.UserUseCase,
) *httpHandler {
	return &httpHandler{
		cfg: cfg,
		uc: &useCase{
			accrual: accrualUseCase,
			balance: balanceUseCase,
			order:   orderUseCase,
			user:    userUseCase,
		},
	}
}
