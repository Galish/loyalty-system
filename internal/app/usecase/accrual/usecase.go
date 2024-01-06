package accrual

import (
	"context"
	"time"

	repo "github.com/Galish/loyalty-system/internal/app/adapters/repository"
	"github.com/Galish/loyalty-system/internal/app/adapters/webapi"
	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/config"
)

const (
	maxAttempts uint = 10
)

type accrualUseCase struct {
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
) *accrualUseCase {
	service := &accrualUseCase{
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

func (uc *accrualUseCase) GetAccrual(ctx context.Context, order *entity.Order) {
	uc.requestCh <- &request{
		order:    string(order.ID),
		user:     order.User,
		attempts: 0,
	}
}

func (uc *accrualUseCase) flushMessages() {
	limiter := newLimiter(uc.limiterInterval)

	for req := range uc.requestCh {
		<-limiter.C

		go func(req *request) {
			accrual, err := uc.accrualAPI.GetAccrual(context.Background(), req.order)
			if err != nil || !accrual.Status.IsFinal() && !req.isAttemptsExceeded() {
				go uc.retry(req)
				return
			}
			accrual.User = req.user

			uc.applyAccrual(context.Background(), accrual)
		}(req)
	}

	limiter.Close()
}

func (uc *accrualUseCase) Close() {
	close(uc.requestCh)
}
