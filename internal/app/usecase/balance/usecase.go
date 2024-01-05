package balance

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type balanceUseCase struct {
	repo repository.BalanceRepository
}

func New(repo repository.BalanceRepository) *balanceUseCase {
	return &balanceUseCase{
		repo: repo,
	}
}
