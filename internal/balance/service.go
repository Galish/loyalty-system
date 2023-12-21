package balance

import (
	"github.com/Galish/loyalty-system/internal/repository"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

type BalanceService struct {
	repo repository.BalanceRepository
}

func NewService(repo repo.BalanceRepository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}
