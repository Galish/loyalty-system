package router

import (
	"loyalty-system/internal/config"
	"net/http"
)

type httpHandler struct {
	cfg *config.Config
}

func newHandler(cfg *config.Config) *httpHandler {
	return &httpHandler{
		cfg: cfg,
	}
}

func (f *httpHandler) stub(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
