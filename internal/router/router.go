package router

import (
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func New(cfg *config.Config, authService *auth.AuthService) *chi.Mux {
	handler := newHandler(cfg, authService)
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithRequestLogger)

		r.Post("/api/user/register", handler.Register)
		r.Post("/api/user/login", handler.Login)

		r.Get("/ping", handler.Ping)
	})

	router.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.WithAuthChecker)
		r.Use(middleware.WithRequestLogger)

		r.Post("/orders", handler.Ping)
		r.Get("/orders", handler.Ping)

		r.Get("/balance", handler.Ping)
		r.Post("/balance/withdraw", handler.Ping)
		r.Get("/withdrawals", handler.Ping)
	})

	return router
}
