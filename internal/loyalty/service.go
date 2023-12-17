package loyalty

import (
	"time"

	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/repository"
)

const limiterInterval = 1 * time.Second

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

func (s *LoyaltyService) Close() {
	close(s.orderCh)
}

func (s *LoyaltyService) flushMessages() {
	limiter := NewLimiter(limiterInterval)

	for order := range s.orderCh {
		<-limiter.C
		s.getOrderAccrual(order)
	}

	limiter.Close()
}

type Limiter struct {
	C        chan struct{}
	ticker   *time.Ticker
	interval time.Duration
}

func NewLimiter(interval time.Duration) *Limiter {
	limiter := Limiter{
		interval: interval,
		C:        make(chan struct{}),
	}

	go func() {
		limiter.C <- struct{}{}

		limiter.ticker = time.NewTicker(limiter.interval)

		for range limiter.ticker.C {
			limiter.C <- struct{}{}
		}
	}()

	return &limiter
}

func (l *Limiter) Close() {
	l.ticker.Stop()
}
