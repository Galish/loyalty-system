package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/loyalty"
	"github.com/Galish/loyalty-system/internal/model"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

type responseBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type requestWithdrawal struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type responseWithdrawal struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (h *httpHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get(auth.AuthHeaderName)

	balance, err := h.loyaltyService.GetBalance(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responseBalance{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.WithError(err).Debug("cannot encode request JSON body")
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		return
	}
}

func (h *httpHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req requestWithdrawal
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		http.Error(w, "cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	withdrawal := model.Withdrawal{
		Order: model.OrderNumber(req.Order),
		Sum:   req.Sum,
		User:  r.Header.Get(auth.AuthHeaderName),
	}

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

func (h *httpHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get(auth.AuthHeaderName)
	withdrawals, err := h.loyaltyService.Withdrawals(r.Context(), user)
	if err != nil {
		logger.WithError(err).Debug("unable to get withdrawals")
		http.Error(w, "unable to get withdrawals", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var res []*responseWithdrawal

	for _, w := range withdrawals {
		res = append(
			res,
			&responseWithdrawal{
				Order:       w.Order.String(),
				Sum:         w.Sum,
				ProcessedAt: w.ProcessedAt.Format(timeLayout),
			},
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.WithError(err).Debug("cannot encode request JSON body")
		http.Error(w, "cannot encode request JSON body", http.StatusInternalServerError)
		return
	}
}
