package psql

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
	repo "github.com/Galish/loyalty-system/internal/app/repository"
)

func (s *psqlStore) CreateOrder(ctx context.Context, order *entity.Order) error {
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
		order.UploadedAt.Value(),
		order.User,
	)

	var (
		uploadedAt time.Time
		user       string
	)

	if err := row.Scan(&uploadedAt, &user); err != nil {
		return err
	}

	if order.UploadedAt.Value().Equal(uploadedAt) {
		return nil
	}

	if order.User != user {
		return repo.ErrOrderConflict
	}

	return repo.ErrOrderExists
}
