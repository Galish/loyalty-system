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

	_, err = tx.ExecContext(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS users (
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
		CREATE UNIQUE INDEX IF NOT EXISTS login_idx ON users (
			login
		)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
