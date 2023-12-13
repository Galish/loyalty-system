package psql

import (
	"context"

	"github.com/Galish/loyalty-system/internal/logger"
)

func (s *psqlStore) Bootstrap(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	logger.Info("database initialization")

	// users
	_, err = tx.ExecContext(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS users (
				_id SERIAL PRIMARY KEY,
				uuid varchar(36) NOT NULL,
				login varchar(250) NOT NULL,
				password varchar(250) NOT NULL,
				is_active boolean DEFAULT true
			)
		`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS login_idx ON users (login)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	// orders
	_, err = tx.ExecContext(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS orders (
				_id SERIAL PRIMARY KEY,
				uuid VARCHAR(250) NOT NULL,
				status VARCHAR(25) NOT NULL,
				accrual INTEGER DEFAULT 0,
				uploaded_at TIMESTAMPTZ NOT NULL,
				user_id VARCHAR(36) NOT NULL
			)
		`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS uuid_idx ON orders (uuid)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	// balance

	_, err = tx.ExecContext(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS balance (
				_id SERIAL PRIMARY KEY,
				user_id VARCHAR(36) NOT NULL,
				points INTEGER DEFAULT 0,
				updated_at TIMESTAMPTZ NOT NULL
			)
		`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS user_id_idx ON balance (user_id)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
