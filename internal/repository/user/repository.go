package userRepository

import (
	"context"
	"database/sql"
	"errors"

	"loyalty-system/internal/config"
	"loyalty-system/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type psqlStore struct {
	db *sql.DB
}

func New(cfg *config.Config) (*psqlStore, error) {
	if cfg.DBAddr == "" {
		return nil, errors.New("database address missing")
	}

	logger.Info("database connection")

	db, err := sql.Open("pgx", cfg.DBAddr)
	if err != nil {
		return nil, err
	}

	store := psqlStore{db}

	if err := store.Bootstrap(context.Background()); err != nil {
		return nil, err
	}

	return &store, nil
}
