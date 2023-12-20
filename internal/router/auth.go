package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug("unable to read request body")
		http.Error(w, "unable to read request body", http.StatusInternalServerError)
		return
	}

	var req authRequest
	if err := json.Unmarshal(body, &req); err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		http.Error(w, "cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		logger.Debug("missing login or password")
		http.Error(w, "missing login or password", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Register(r.Context(), req.Login, req.Password)
	if errors.Is(err, repo.ErrUserConflict) {
		logger.WithError(err).Debug("unable to write to repository")
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		logger.WithError(err).Debug("unable to write to repository")
		http.Error(w, "unable to write to repository", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token)

	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug("unable to read request body")
		http.Error(w, "unable to read request body", http.StatusInternalServerError)
		return
	}

	var req authRequest
	if err := json.Unmarshal(body, &req); err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		http.Error(w, "cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		logger.Debug("missing login or password")
		http.Error(w, "missing login or password", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(r.Context(), req.Login, req.Password)
	if errors.Is(err, auth.ErrIncorrectLoginPassword) || errors.Is(err, repo.ErrUserNotFound) {
		logger.WithError(err).Debug("unable to authenticate")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		logger.WithError(err).Debug("unable to authenticate")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token)

	w.WriteHeader(http.StatusOK)
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:  auth.AuthCookieName,
			Value: token,
			Path:  "/",
		},
	)
}
