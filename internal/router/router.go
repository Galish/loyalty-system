package router

import (
	"loyalty-system/internal/config"

	"github.com/go-chi/chi/v5"
)

func New(cfg *config.Config) *chi.Mux {
	handler := newHandler(cfg)
	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handler.stub)
		r.Post("/login", handler.stub)

		r.Post("/orders", handler.stub)
		r.Get("/orders", handler.stub)

		r.Get("/balance", handler.stub)
		r.Post("/balance/withdraw", handler.stub)
		r.Get("/withdrawals", handler.stub)
	})

	return router
}
