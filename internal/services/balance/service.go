package balance

import (
	"context"

	"github.com/Galish/loyalty-system/internal/entity"
	"github.com/Galish/loyalty-system/internal/repository"
)

type BalanceManager interface {
	GetBalance(context.Context, string) (*entity.Balance, error)
	Withdraw(context.Context, *entity.Withdrawal) error
	Withdrawals(context.Context, string) ([]*entity.Withdrawal, error)
}

type BalanceService struct {
	repo repository.BalanceRepository
}

func NewService(repo repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}
