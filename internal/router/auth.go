package router

import (
	"encoding/json"
	"errors"
	"io"

	userRepo "loyalty-system/internal/repository/user"
	"net/http"
)

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *httpHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user credentials
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.authService.Authenticate(r.Context(), user.Login, user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// TODO: set cookie

	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user credentials
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.authService.Register(r.Context(), user.Login, user.Password)
	if err != nil && errors.Is(err, userRepo.ErrConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
