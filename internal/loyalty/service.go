package loyalty

import "github.com/Galish/loyalty-system/internal/repository"

type LoyaltyService struct {
	repo repository.LoyaltyRepository
}

func NewService(repo repository.LoyaltyRepository) *LoyaltyService {
	return &LoyaltyService{repo}
}
