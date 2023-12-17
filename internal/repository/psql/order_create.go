package psql

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/logger"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *psqlStore) CreateOrder(ctx context.Context, order *repo.Order) error {
	row := s.db.QueryRowContext(
		ctx,
		`
			INSERT INTO orders (uuid, status, accrual, uploaded_at, user_id)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (uuid)
			DO UPDATE SET uuid=excluded.uuid
			RETURNING uploaded_at, user_id
		`,
		order.ID,
		order.Status,
		order.Accrual,
		order.UploadedAt,
		order.User,
	)

	var (
		uploadedAt time.Time
		user       string
	)

	if err := row.Scan(&uploadedAt, &user); err != nil {
		return err
	}

	logger.WithFields(logger.Fields{
		"Order":      order,
		"UploadedAt": uploadedAt,
		"User":       user,
	}).Debug("creating order error")

	if !order.UploadedAt.Equal(uploadedAt) {
		if order.User != user {
			return repo.ErrOrderConflict
		}

		return repo.ErrOrderExists
	}

	return nil
}
