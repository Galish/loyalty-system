package loyalty

import (
	"time"

	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/repository"
)

type LoyaltyService struct {
	repo    repository.LoyaltyRepository
	cfg     *config.Config
	orderCh chan *Order
}

func NewService(repo repository.LoyaltyRepository, cfg *config.Config) *LoyaltyService {
	service := LoyaltyService{
		repo:    repo,
		cfg:     cfg,
		orderCh: make(chan *Order),
	}

	go service.flushMessages()

	return &service
}

func (s *LoyaltyService) flushMessages() {
	limiter := NewLimiter(10 * time.Second)

	for order := range s.orderCh {
		<-limiter
		s.getOrderAccrual(order)
	}
}

func NewLimiter(duration time.Duration) chan struct{} {
	ch := make(chan struct{})

	go func() {
		ch <- struct{}{}

		ticker := time.NewTicker(duration)

		for range ticker.C {
			ch <- struct{}{}
		}
	}()

	return ch
}
