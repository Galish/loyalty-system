package psql

import (
	"context"
	"fmt"
	"time"
)

func (s *psqlStore) UpdateBalance(ctx context.Context, user string, value int) error {
	fmt.Println("Balance", user, value)

	_, err := s.db.ExecContext(
		ctx,
		`
			INSERT INTO balance (user_id, points, updated_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id)
			DO UPDATE SET points = balance.points + excluded.points
		`,
		user,
		value,
		time.Now(),
	)
	// , updated_at = excluded.updated_at
	if err != nil {
		fmt.Println("Err", err)
		return err
	}

	return nil
}
