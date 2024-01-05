package order

import (
	"context"
	"errors"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

var ErrInvalidOrderNumber = errors.New("invalid order number value")

func (uc *orderUseCase) AddOrder(ctx context.Context, order entity.Order) error {
	if !order.IsValid() {
		return ErrInvalidOrderNumber
	}

	order.Status = entity.StatusNew
	order.UploadedAt = entity.Time(time.Now())

	if err := uc.repo.CreateOrder(ctx, &order); err != nil {
		return err
	}

	return nil
}
