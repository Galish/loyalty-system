package psql

import (
	"context"
	"fmt"
	"time"
)

func (s *psqlStore) UpdateBalance(ctx context.Context, user string, value float32) error {
	fmt.Println("Balance", user, value)

	_, err := s.db.ExecContext(
		ctx,
		`
			INSERT INTO balance (user_id, current, updated_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id)
			DO UPDATE
			SET current = balance.current + excluded.current
		`,
		user,
		value,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
