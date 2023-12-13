package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/loyalty"
)

func (h *httpHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get(auth.AuthHeaderName)

	balance, err := h.loyaltyService.GetBalance(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	var withdraw loyalty.Withdraw
	err := json.NewDecoder(r.Body).Decode(&withdraw)
	if err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	withdraw.User = r.Header.Get(auth.AuthHeaderName)

	err = h.loyaltyService.Withdraw(r.Context(), &withdraw)
	if err != nil {
		var httpStatus = http.StatusInternalServerError

		if errors.Is(err, loyalty.ErrInsufficientFunds) {
			httpStatus = http.StatusPaymentRequired
		}

		if errors.Is(err, loyalty.ErrInvalidOrderNumber) {
			httpStatus = http.StatusUnprocessableEntity
		}

		logger.WithError(err).Debug("unable to withdraw funds")
		http.Error(w, "unable to withdraw funds", httpStatus)
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
