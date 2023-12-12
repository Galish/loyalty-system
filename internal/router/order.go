package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/loyalty"
	"github.com/Galish/loyalty-system/internal/repository"
)

func (h *httpHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	number, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug("unable to read request body")
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get(auth.AuthHeaderName)

	_, err = h.loyaltyService.AddOrder(r.Context(), string(number), userID)
	if err != nil {
		if errors.Is(loyalty.ErrInvalidOrderNumber, err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(repository.ErrOrderConflict, err) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if errors.Is(repository.ErrOrderExists, err) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		}

		logger.WithError(err).Debug("unable to add order")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *httpHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(auth.AuthHeaderName)

	orders, err := h.loyaltyService.GetOrders(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Debug("unable to read from repository")
		http.Error(w, "unable to read from repository", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		logger.WithError(err).Debug("cannot encode request JSON body")
		http.Error(w, "cannot encode request JSON body", http.StatusBadRequest)
		return
	}
}
