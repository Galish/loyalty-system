package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Galish/loyalty-system/internal/app/entity"
	repo "github.com/Galish/loyalty-system/internal/app/repository"
	"github.com/Galish/loyalty-system/internal/app/validation"
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
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

	balance, err := h.balanceService.GetBalance(r.Context(), user)
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
		logger.WithError(err).Debug(errEncodeResponseBody)
		http.Error(w, errEncodeResponseBody, http.StatusInternalServerError)
		return
	}
}

func (h *httpHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req requestWithdrawal
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.WithError(err).Debug(errDecodeRequestBody)
		http.Error(w, errDecodeRequestBody, http.StatusBadRequest)
		return
	}

	withdrawal := entity.Withdrawal{
		Order: req.Order,
		Sum:   req.Sum,
		User:  r.Header.Get(auth.AuthHeaderName),
	}

	err = h.balanceService.Withdraw(r.Context(), &withdrawal)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	logger.WithError(err).Debug(errWithdrawFunds)

	if errors.Is(err, validation.ErrInvalidOrderNumber) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if errors.Is(err, repo.ErrInsufficientFunds) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	http.Error(w, errWithdrawFunds, http.StatusInternalServerError)
}

func (h *httpHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	withdrawals, err := h.balanceService.Withdrawals(
		r.Context(),
		r.Header.Get(auth.AuthHeaderName),
	)
	if err != nil {
		logger.WithError(err).Debug(errGetWithdrawals)
		http.Error(w, errGetWithdrawals, http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var res []*responseWithdrawal
	for _, withdrawal := range withdrawals {
		res = append(
			res,
			&responseWithdrawal{
				Order:       withdrawal.Order,
				Sum:         withdrawal.Sum,
				ProcessedAt: withdrawal.ProcessedAt.Format(),
			},
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.WithError(err).Debug(errEncodeResponseBody)
		http.Error(w, errEncodeResponseBody, http.StatusInternalServerError)
		return
	}
}
