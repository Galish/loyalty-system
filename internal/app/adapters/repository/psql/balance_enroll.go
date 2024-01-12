package psql

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (s *psqlStore) Enroll(ctx context.Context, enroll *entity.Enrollment) error {
	_, err := s.db.ExecContext(
		ctx,
		`
			INSERT INTO balance (user_id, current, updated_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id)
			DO UPDATE
			SET current = balance.current + excluded.current
		`,
		enroll.User,
		enroll.Sum,
		enroll.ProcessedAt,
	)
	if err != nil {
		return err
	}

	return nil
}