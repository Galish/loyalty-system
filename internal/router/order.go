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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug("unable to read request body")
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}

	newOrder := loyalty.Order{
		ID:   loyalty.OrderNumber(string(body)),
		User: r.Header.Get(auth.AuthHeaderName),
	}

	err = h.loyaltyService.AddOrder(r.Context(), &newOrder)
	if err != nil {
		logger.WithError(err).Debug("unable to add order")

		if errors.Is(err, loyalty.ErrIncorrectOrderNumber) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, repository.ErrOrderConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		if errors.Is(err, repository.ErrOrderExists) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(err.Error()))
			return
		}

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
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		return
	}
}
