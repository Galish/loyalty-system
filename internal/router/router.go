package router

import (
	"loyalty-system/internal/auth"
	"loyalty-system/internal/config"

	"github.com/go-chi/chi/v5"
)

func New(cfg *config.Config, authService *auth.AuthService) *chi.Mux {
	handler := newHandler(cfg, authService)
	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		r.Post("/orders", handler.stub)
		r.Get("/orders", handler.stub)

		r.Get("/balance", handler.stub)
		r.Post("/balance/withdraw", handler.stub)
		r.Get("/withdrawals", handler.stub)
	})

	return router
}
