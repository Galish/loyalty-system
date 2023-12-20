package psql

import (
	"context"
	"errors"

	"github.com/Galish/loyalty-system/internal/model"
	repo "github.com/Galish/loyalty-system/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *psqlStore) Withdraw(ctx context.Context, withdrawal *model.Withdrawal) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		`
			UPDATE balance
			SET current = balance.current - $1, withdrawn = balance.withdrawn + $1
			WHERE user_id = $2
		`,
		withdrawal.Sum,
		withdrawal.User,
	)
	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return repo.ErrInsufficientFunds
		}

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()

		return err
	}

	if rows == 0 {
		tx.Rollback()

		return repo.ErrInsufficientFunds
	}

	_, err = tx.Exec(
		`
			INSERT INTO withdrawals (order_id, user_id, sum, processed_at)
			VALUES ($1, $2, $3, $4)
		`,
		withdrawal.Order,
		withdrawal.User,
		withdrawal.Sum,
		withdrawal.ProcessedAt,
	)

	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}
