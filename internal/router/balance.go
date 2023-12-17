package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/loyalty"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (h *httpHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get(auth.AuthHeaderName)

	balance, err := h.loyaltyService.GetBalance(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(balance); err != nil {
		logger.WithError(err).Debug("cannot encode request JSON body")
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		return
	}
}

func (h *httpHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var withdrawal loyalty.Withdrawal
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	if err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		http.Error(w, "cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	withdrawal.User = r.Header.Get(auth.AuthHeaderName)

	err = h.loyaltyService.Withdraw(r.Context(), &withdrawal)
	if err != nil {
		logger.WithError(err).Debug("unable to withdraw funds")

		if errors.Is(err, repo.ErrInsufficientFunds) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}

		if errors.Is(err, loyalty.ErrIncorrectOrderNumber) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		http.Error(w, "unable to withdraw funds", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get(auth.AuthHeaderName)
	withdrawals, err := h.loyaltyService.GetWithdrawals(r.Context(), user)
	if err != nil {
		logger.WithError(err).Debug("unable to get withdrawals")
		http.Error(w, "unable to get withdrawals", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
		logger.WithError(err).Debug("cannot encode request JSON body")
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		return
	}
}
