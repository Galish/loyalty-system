package order

import (
	"context"
	"errors"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/datetime"
)

var ErrInvalidOrderNumber = errors.New("invalid order number value")

func (uc *orderUseCase) AddOrder(ctx context.Context, order entity.Order) error {
	if !order.IsValid() {
		return ErrInvalidOrderNumber
	}

	order.Status = entity.StatusNew
	order.UploadedAt = datetime.Round(time.Now())

	if err := uc.repo.CreateOrder(ctx, &order); err != nil {
		return err
	}

	return nil
}
