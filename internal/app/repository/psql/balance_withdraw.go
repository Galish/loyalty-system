package psql

import (
	"context"
	"errors"

	"github.com/Galish/loyalty-system/internal/app/entity"
	repo "github.com/Galish/loyalty-system/internal/app/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *psqlStore) Withdraw(ctx context.Context, withdrawal *entity.Withdrawal) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		`
			UPDATE balance
			SET current = balance.current - $1,
				withdrawn = balance.withdrawn + $1,
				updated_at = $2
			WHERE user_id = $3
		`,
		withdrawal.Sum,
		withdrawal.ProcessedAt,
		withdrawal.User,
	)
	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == errCodeCheckViolation {
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
