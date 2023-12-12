package main

import (
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/loyalty"
	"github.com/Galish/loyalty-system/internal/repository/psql"
	"github.com/Galish/loyalty-system/internal/router"
	"github.com/Galish/loyalty-system/internal/server"
)

func main() {
	cfg := config.New()

	logger.Initialize(cfg.LogLevel)

	store, err := psql.NewStore(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	authService := auth.NewService(store, cfg.SecretKey)
	loyaltyService := loyalty.NewService(store)

	router := router.New(cfg, authService, loyaltyService)

	httpServer := server.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
