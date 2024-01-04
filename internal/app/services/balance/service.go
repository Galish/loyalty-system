package balance

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type BalanceManager interface {
	GetBalance(context.Context, string) (*entity.Balance, error)
	Withdraw(context.Context, *entity.Withdrawal) error
	Withdrawals(context.Context, string) ([]*entity.Withdrawal, error)
}

type balanceService struct {
	repo repository.BalanceRepository
}

func New(repo repository.BalanceRepository) *balanceService {
	return &balanceService{
		repo: repo,
	}
}
