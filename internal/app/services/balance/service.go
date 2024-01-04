package balance

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type balanceService struct {
	repo repository.BalanceRepository
}

func New(repo repository.BalanceRepository) *balanceService {
	return &balanceService{
		repo: repo,
	}
}
