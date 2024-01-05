package main

import (
	"github.com/Galish/loyalty-system/internal/app/handlers"
	"github.com/Galish/loyalty-system/internal/app/repository/psql"
	"github.com/Galish/loyalty-system/internal/app/usecase/accrual"
	"github.com/Galish/loyalty-system/internal/app/usecase/balance"
	"github.com/Galish/loyalty-system/internal/app/usecase/order"
	"github.com/Galish/loyalty-system/internal/app/usecase/user"
	"github.com/Galish/loyalty-system/internal/app/webapi"
	"github.com/Galish/loyalty-system/internal/config"
	httpserver "github.com/Galish/loyalty-system/internal/http/server"
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

	webAPI := webapi.New(cfg)

	accrualUseCase := accrual.New(webAPI, store, store, cfg)
	defer accrualUseCase.Close()

	balanceUseCase := balance.New(store)
	orderUseCase := order.New(store)
	userUseCase := user.New(store, cfg.SecretKey)

	handler := handlers.NewHandler(
		cfg,
		accrualUseCase,
		balanceUseCase,
		orderUseCase,
		userUseCase,
	)

	router := handlers.NewRouter(cfg, handler)

	httpServer := httpserver.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
