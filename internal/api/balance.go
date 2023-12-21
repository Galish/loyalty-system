package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
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

	orderNumber := model.OrderNumber(req.Order)
	if !orderNumber.IsValid() {
		http.Error(w, errInvalidOrderNumber, http.StatusUnprocessableEntity)
		return
	}

	withdrawal := model.Withdrawal{
		Order: orderNumber,
		Sum:   req.Sum,
		User:  r.Header.Get(auth.AuthHeaderName),
	}

	err = h.balanceService.Withdraw(r.Context(), &withdrawal)
	if err != nil {
		logger.WithError(err).Debug(errWithdrawFunds)

		if errors.Is(err, repo.ErrInsufficientFunds) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}

		http.Error(w, errWithdrawFunds, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get(auth.AuthHeaderName)
	withdrawals, err := h.balanceService.Withdrawals(r.Context(), user)
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
		logger.WithError(err).Debug(errEncodeResponseBody)
		http.Error(w, errEncodeResponseBody, http.StatusInternalServerError)
		return
	}
}
