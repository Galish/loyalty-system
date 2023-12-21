package main

import (
	"github.com/Galish/loyalty-system/internal/accrual"
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/balance"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/order"
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
	orderService := order.NewService(store)
	balanceService := balance.NewService(store)
	accrualService := accrual.NewService(store, store, cfg.AccrualAddr)
	defer accrualService.Close()

	router := router.New(
		cfg,
		authService,
		orderService,
		balanceService,
		accrualService,
	)

	httpServer := server.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
