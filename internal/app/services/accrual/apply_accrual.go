package accrual

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/logger"
)

func (s *accrualService) applyAccrual(ctx context.Context, accrual *entity.Accrual) error {
	err := s.orderRepo.UpdateOrder(
		ctx,
		&entity.Order{
			ID:      accrual.Order,
			Status:  accrual.Status,
			Accrual: accrual.Value,
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update order")
		return err
	}

	err = s.balanceRepo.Enroll(
		ctx,
		&entity.Enrollment{
			User:        accrual.User,
			Sum:         accrual.Value,
			ProcessedAt: entity.Time(time.Now()),
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update balance")
		return err
	}

	return nil
}
