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

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug("unable to read request body")
		http.Error(w, "unable to read request body", http.StatusInternalServerError)
		return
	}

	var creds auth.Credentials
	if err := json.Unmarshal(body, &creds); err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		http.Error(w, "cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	if creds.Login == "" || creds.Password == "" {
		logger.Debug("missing login or password")
		http.Error(w, "missing login or password", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Register(r.Context(), creds)
	if err != nil && errors.Is(err, repo.ErrUserConflict) {
		logger.WithError(err).Debug("unable to write to repository")
		http.Error(w, "unable to write to repository", http.StatusConflict)
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

	var creds auth.Credentials
	if err := json.Unmarshal(body, &creds); err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		http.Error(w, "cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	if creds.Login == "" || creds.Password == "" {
		logger.Debug("missing login or password")
		http.Error(w, "missing login or password", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(r.Context(), creds)
	if err != nil {
		logger.WithError(err).Debug("unable to write to repository")
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
