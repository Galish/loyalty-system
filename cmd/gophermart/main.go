package main

import (
	"github.com/Galish/loyalty-system/internal/app/handlers"
	"github.com/Galish/loyalty-system/internal/app/repository/psql"
	"github.com/Galish/loyalty-system/internal/app/services/accrual"
	"github.com/Galish/loyalty-system/internal/app/services/auth"
	"github.com/Galish/loyalty-system/internal/app/services/balance"
	"github.com/Galish/loyalty-system/internal/app/services/order"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/http/httpserver"
	"github.com/Galish/loyalty-system/internal/logger"
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
		handlers.NewHandler(
			cfg,
			authService,
			orderService,
			balanceService,
			accrualService,
		),
		authService,
	)

	httpServer := httpserver.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
