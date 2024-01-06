package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	repo "github.com/Galish/loyalty-system/internal/app/adapters/repository"
	"github.com/Galish/loyalty-system/internal/app/usecase/user"
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
)

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug(errReadRequestBody)
		http.Error(w, errReadRequestBody, http.StatusInternalServerError)
		return
	}

	var req authRequest
	if err := json.Unmarshal(body, &req); err != nil {
		logger.WithError(err).Debug(errDecodeRequestBody)
		http.Error(w, errDecodeRequestBody, http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		logger.Debug(errMissingLoginOrPassword)
		http.Error(w, errMissingLoginOrPassword, http.StatusBadRequest)
		return
	}

	token, err := h.uc.user.Register(r.Context(), req.Login, req.Password)
	if errors.Is(err, repo.ErrUserConflict) {
		logger.WithError(err).Debug(errRegisterUser)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		logger.WithError(err).Debug(errRegisterUser)
		http.Error(w, errRegisterUser, http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token)

	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug(errReadRequestBody)
		http.Error(w, errReadRequestBody, http.StatusInternalServerError)
		return
	}

	var req authRequest
	if err := json.Unmarshal(body, &req); err != nil {
		logger.WithError(err).Debug(errDecodeRequestBody)
		http.Error(w, errDecodeRequestBody, http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		logger.Debug(errMissingLoginOrPassword)
		http.Error(w, errMissingLoginOrPassword, http.StatusBadRequest)
		return
	}

	token, err := h.uc.user.Authenticate(r.Context(), req.Login, req.Password)
	if errors.Is(err, user.ErrIncorrectLoginPassword) || errors.Is(err, repo.ErrUserNotFound) {
		logger.WithError(err).Debug(errAuthenticate)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		logger.WithError(err).Debug(errAuthenticate)
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
