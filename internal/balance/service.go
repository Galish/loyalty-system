package balance

import (
	"github.com/Galish/loyalty-system/internal/repository"
)

type BalanceService struct {
	repo repository.BalanceRepository
}

func NewService(repo repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}
