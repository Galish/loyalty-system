package main

import (
	"loyalty-system/internal/config"
	"loyalty-system/internal/router"
	"loyalty-system/internal/server"
)

func main() {
	cfg := config.New()
	router := router.New(cfg)

	httpServer := server.New(cfg.SrvAddr, router)
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
