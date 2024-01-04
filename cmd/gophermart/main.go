package main

import (
	"github.com/Galish/loyalty-system/internal/app/handlers"
	"github.com/Galish/loyalty-system/internal/app/repository/psql"
	"github.com/Galish/loyalty-system/internal/app/services"
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

	svc := services.New(
		cfg,
		store,
		webapi.New(cfg),
	)
	defer svc.Close()

	handler := handlers.NewHandler(cfg, svc)
	router := handlers.NewRouter(cfg, handler)

	httpServer := httpserver.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
