package balance

import (
	"context"

	"github.com/Galish/loyalty-system/internal/model"
	"github.com/Galish/loyalty-system/internal/repository"
)

type BalanceManager interface {
	GetBalance(context.Context, string) (*model.Balance, error)
	Withdraw(context.Context, *model.Withdrawal) error
	Withdrawals(context.Context, string) ([]*model.Withdrawal, error)
}

type BalanceService struct {
	repo repository.BalanceRepository
}

func NewService(repo repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}
