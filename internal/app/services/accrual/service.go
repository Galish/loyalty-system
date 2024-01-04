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

type accrualService struct {
	accrualAPI      webapi.AccrualGetter
	orderRepo       repo.OrderRepository
	balanceRepo     repo.BalanceRepository
	addr            string
	limiterInterval time.Duration
	requestCh       chan *request
}

func New(
	accrualAPI webapi.AccrualGetter,
	orderRepo repo.OrderRepository,
	balanceRepo repo.BalanceRepository,
	cfg *config.Config,
) *accrualService {
	service := &accrualService{
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

func (s *accrualService) GetAccrual(ctx context.Context, order *entity.Order) {
	s.requestCh <- &request{
		order:    string(order.ID),
		user:     order.User,
		attempts: 0,
	}
}

func (s *accrualService) flushMessages() {
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

func (s *accrualService) Close() {
	close(s.requestCh)
}
