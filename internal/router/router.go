package router

import (
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/loyalty"
	"github.com/Galish/loyalty-system/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func New(
	cfg *config.Config,
	auth *auth.AuthService,
	loyalty *loyalty.LoyaltyService,
) *chi.Mux {
	handler := newHandler(cfg, auth, loyalty)
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

		r.Post("/orders", handler.AddOrder)
		r.Get("/orders", handler.GetOrders)

		r.Get("/balance", handler.Ping)
		r.Post("/balance/withdraw", handler.Ping)
		r.Get("/withdrawals", handler.Ping)
	})

	return router
}
