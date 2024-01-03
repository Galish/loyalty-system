package handlers

import (
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/middleware"
	"github.com/Galish/loyalty-system/internal/services/accrual"
	"github.com/Galish/loyalty-system/internal/services/auth"
	"github.com/Galish/loyalty-system/internal/services/balance"
	"github.com/Galish/loyalty-system/internal/services/order"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	cfg *config.Config,
	auth auth.AuthManager,
	order order.OrderManager,
	balance balance.BalanceManager,
	accrual accrual.AccrualManager,
) *chi.Mux {
	handler := newHandler(cfg, auth, order, balance, accrual)
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

		r.Get("/ping", handler.Health)
	})

	return router
}
