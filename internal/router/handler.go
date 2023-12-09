package router

import (
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/loyalty"
)

type httpHandler struct {
	cfg            *config.Config
	authService    *auth.AuthService
	loyaltyService *loyalty.LoyaltyService
}

func newHandler(
	cfg *config.Config,
	auth *auth.AuthService,
	loyalty *loyalty.LoyaltyService,
) *httpHandler {
	return &httpHandler{
		cfg:            cfg,
		authService:    auth,
		loyaltyService: loyalty,
	}
}

func (h *httpHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
