package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"loyalty-system/internal/auth"
	userRepo "loyalty-system/internal/repository/user"
)

const AuthCookieName = "auth"

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var creds auth.Credentials
	if err := json.Unmarshal(body, &creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Register(r.Context(), creds)
	if err != nil && errors.Is(err, userRepo.ErrConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var creds auth.Credentials
	if err := json.Unmarshal(body, &creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(r.Context(), creds)
	if err != nil {
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
			Name:  AuthCookieName,
			Value: token,
		},
	)
}
