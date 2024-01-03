package main

import (
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/handlers"
	"github.com/Galish/loyalty-system/internal/httpserver"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/repository/psql"
	"github.com/Galish/loyalty-system/internal/services/accrual"
	"github.com/Galish/loyalty-system/internal/services/auth"
	"github.com/Galish/loyalty-system/internal/services/balance"
	"github.com/Galish/loyalty-system/internal/services/order"
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
	accrualService := accrual.NewService(store, store, cfg)
	defer accrualService.Close()

	router := handlers.NewRouter(
		cfg,
		authService,
		orderService,
		balanceService,
		accrualService,
	)

	httpServer := httpserver.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
