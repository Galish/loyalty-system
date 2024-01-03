package accrual

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/entity"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

const (
	maxAttempts uint = 10
)

type AccrualManager interface {
	GetAccrual(context.Context, *entity.Order)
}

type AccrualService struct {
	orderRepo       repo.OrderRepository
	balanceRepo     repo.BalanceRepository
	addr            string
	limiterInterval time.Duration
	requestCh       chan *request
}

func NewService(
	orderRepo repo.OrderRepository,
	balanceRepo repo.BalanceRepository,
	cfg *config.Config,
) *AccrualService {
	service := &AccrualService{
		orderRepo:       orderRepo,
		balanceRepo:     balanceRepo,
		addr:            cfg.AccrualAddr,
		limiterInterval: time.Duration(cfg.AccrualInterval) * time.Millisecond,
		requestCh:       make(chan *request),
	}

	go service.flushMessages()

	return service
}

func (s *AccrualService) GetAccrual(ctx context.Context, order *entity.Order) {
	s.requestCh <- &request{
		order:    string(order.ID),
		user:     order.User,
		attempts: 0,
	}
}

func (s *AccrualService) flushMessages() {
	limiter := newLimiter(s.limiterInterval)

	for req := range s.requestCh {
		<-limiter.C

		go func(req *request) {
			accrual, err := s.fetchAccrual(context.Background(), req)
			if err != nil || !accrual.Status.IsFinal() && !req.isAttemptsExceeded() {
				go s.retry(req)
				return
			}

			s.applyAccrual(context.Background(), accrual)
		}(req)
	}

	limiter.Close()
}

func (s *AccrualService) Close() {
	close(s.requestCh)
}
