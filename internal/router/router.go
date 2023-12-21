package router

import (
	"github.com/Galish/loyalty-system/internal/accrual"
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
	accrual accrual.AccrualManager,
) *chi.Mux {
	handler := newHandler(cfg, auth, loyalty, accrual)
	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.WithRequestLogger)

			r.Post("/register", handler.Register)
			r.Post("/login", handler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithAuthChecker(auth))
			r.Use(middleware.WithRequestLogger)

			r.Post("/orders", handler.AddOrder)
			r.Get("/orders", handler.GetOrders)

			r.Get("/balance", handler.GetBalance)
			r.Post("/balance/withdraw", handler.Withdraw)
			r.Get("/withdrawals", handler.Withdrawals)
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.WithRequestLogger)

		r.Get("/ping", handler.Ping)
	})

	return router
}
