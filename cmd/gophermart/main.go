package main

import (
	"loyalty-system/internal/auth"
	"loyalty-system/internal/config"
	"loyalty-system/internal/logger"
	"loyalty-system/internal/repository/psql"
	"loyalty-system/internal/router"
	"loyalty-system/internal/server"
)

func main() {
	cfg := config.New()

	logger.Initialize(cfg.LogLevel)

	store, err := psql.NewStore(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	authService := auth.NewService(store)
	router := router.New(cfg, authService)

	httpServer := server.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
