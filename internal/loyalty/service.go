package loyalty

import (
	"time"

	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/repository"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

const limiterInterval = 1 * time.Second

type LoyaltyService struct {
	repo repository.LoyaltyRepository
	cfg  *config.Config
}

func NewService(repo repo.LoyaltyRepository, cfg *config.Config) *LoyaltyService {
	return &LoyaltyService{
		repo: repo,
		cfg:  cfg,
	}
}
