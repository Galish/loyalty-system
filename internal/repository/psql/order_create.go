package psql

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/repository"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *psqlStore) CreateOrder(ctx context.Context, order *repository.Order) error {
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

	if !order.UploadedAt.Equal(uploadedAt) && order.User != user {
		return repo.ErrOrderConflict
	}

	if !order.UploadedAt.Equal(uploadedAt) {
		return repo.ErrOrderExists
	}

	return nil
}
