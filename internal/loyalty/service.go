package loyalty

import (
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/repository"
)

type LoyaltyService struct {
	repo    repository.LoyaltyRepository
	cfg     *config.Config
	orderCh chan *Order
	accrual *accrualClient
}

func NewService(repo repository.LoyaltyRepository, cfg *config.Config) *LoyaltyService {
	service := &LoyaltyService{
		repo:    repo,
		cfg:     cfg,
		orderCh: make(chan *Order),
		accrual: newAccrualClient(repo, cfg),
	}

	go service.flushMessages()

	return service
}

func (s *LoyaltyService) flushMessages() {
	for order := range s.orderCh {
		go s.accrual.newOrder(order)
	}
}

func (s *LoyaltyService) Close() {
	close(s.orderCh)
	s.accrual.close()
}
