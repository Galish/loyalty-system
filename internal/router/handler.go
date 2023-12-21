package router

import (
	"net/http"

	"github.com/Galish/loyalty-system/internal/accrual"
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/loyalty"
)

type httpHandler struct {
	cfg            *config.Config
	authService    *auth.AuthService
	loyaltyService *loyalty.LoyaltyService
	accrualService accrual.AccrualManager
}

func newHandler(
	cfg *config.Config,
	auth *auth.AuthService,
	loyalty *loyalty.LoyaltyService,
	accrual accrual.AccrualManager,
) *httpHandler {
	return &httpHandler{
		cfg:            cfg,
		authService:    auth,
		loyaltyService: loyalty,
		accrualService: accrual,
	}
}

func (h *httpHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
