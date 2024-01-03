package accrual

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
	repo "github.com/Galish/loyalty-system/internal/app/repository"
	"github.com/Galish/loyalty-system/internal/app/webapi"
	"github.com/Galish/loyalty-system/internal/config"
)

const (
	maxAttempts uint = 10
)

type AccrualManager interface {
	GetAccrual(context.Context, *entity.Order)
}

type AccrualService struct {
	accrualAPI      webapi.AccrualGetter
	orderRepo       repo.OrderRepository
	balanceRepo     repo.BalanceRepository
	addr            string
	limiterInterval time.Duration
	requestCh       chan *request
}

func NewService(
	accrualAPI webapi.AccrualGetter,
	orderRepo repo.OrderRepository,
	balanceRepo repo.BalanceRepository,
	cfg *config.Config,
) *AccrualService {
	service := &AccrualService{
		accrualAPI:      accrualAPI,
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
			accrual, err := s.accrualAPI.GetAccrual(context.Background(), req.order)
			if err != nil || !accrual.Status.IsFinal() && !req.isAttemptsExceeded() {
				go s.retry(req)
				return
			}
			accrual.User = req.user

			s.applyAccrual(context.Background(), accrual)
		}(req)
	}

	limiter.Close()
}

func (s *AccrualService) Close() {
	close(s.requestCh)
}
