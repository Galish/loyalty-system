package router

import (
	"loyalty-system/internal/auth"
	"loyalty-system/internal/config"
	"net/http"
)

type httpHandler struct {
	cfg         *config.Config
	authService *auth.AuthService
}

func newHandler(cfg *config.Config, authService *auth.AuthService) *httpHandler {
	return &httpHandler{
		cfg:         cfg,
		authService: authService,
	}
}

func (f *httpHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
