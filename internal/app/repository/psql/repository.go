package psql

import (
	"database/sql"
	"errors"
	"os"

	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	errCodeConflict       = "23505"
	errCodeCheckViolation = "23514"
)

type psqlStore struct {
	db *sql.DB
}

func NewStore(cfg *config.Config) (*psqlStore, error) {
	if cfg.DBAddr == "" {
		return nil, errors.New("database address missing")
	}

	logger.Info("database connection")

	db, err := sql.Open("pgx", cfg.DBAddr)
	if err != nil {
		return nil, err
	}

	store := psqlStore{db}

	logger.Info("database initialization")

	if err := store.init(); err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *psqlStore) init() error {
	query, err := os.ReadFile("internal/app/repository/psql/init.sql")
	if err != nil {
		return err
	}

	_, err = s.db.Exec(string(query))
	if err != nil {
		return err
	}

	return nil
}

func (s *psqlStore) Close() error {
	return s.db.Close()
}
