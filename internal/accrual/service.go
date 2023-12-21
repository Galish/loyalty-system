package accrual

import (
	"time"

	"github.com/Galish/loyalty-system/internal/model"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

const limiterInterval time.Duration = 1 * time.Second

type AccrualManager interface {
	GetAccrual(order *model.Order)
}

type AccrualService struct {
	orderRepo   repo.OrderRepository
	balanceRepo repo.BalanceRepository
	addr        string
	requestCh   chan *request
}

func NewService(
	orderRepo repo.OrderRepository,
	balanceRepo repo.BalanceRepository,
	addr string,
) *AccrualService {
	service := &AccrualService{
		orderRepo:   orderRepo,
		balanceRepo: balanceRepo,
		addr:        addr,
		requestCh:   make(chan *request),
	}

	go service.flushMessages()

	return service
}

func (s *AccrualService) GetAccrual(order *model.Order) {
	s.requestCh <- &request{
		order:    string(order.ID),
		user:     order.User,
		attempts: 0,
	}
}

func (s *AccrualService) flushMessages() {
	limiter := newLimiter(limiterInterval)

	for req := range s.requestCh {
		<-limiter.C

		go func(req *request) {
			accrual, err := s.fetchAccrual(req)
			if err != nil || !accrual.Status.IsFinal() && !req.isAttemptsExceeded() {
				go s.retry(req)
				return
			}

			s.applyAccrual(accrual)
		}(req)
	}

	limiter.Close()
}

func (s *AccrualService) Close() {
	close(s.requestCh)
}
