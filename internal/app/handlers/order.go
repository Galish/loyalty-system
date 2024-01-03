package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/app/repository"
	"github.com/Galish/loyalty-system/internal/app/services/auth"
	"github.com/Galish/loyalty-system/internal/logger"
)

type orderResponse struct {
	ID         string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

func (h *httpHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug(errReadRequestBody)
		http.Error(w, errReadRequestBody, http.StatusBadRequest)
		return
	}

	newOrder := entity.Order{
		ID:   entity.OrderNumber(string(body)),
		User: r.Header.Get(auth.AuthHeaderName),
	}

	err = h.orderService.AddOrder(r.Context(), newOrder)
	if err == nil {
		go h.accrualService.GetAccrual(context.Background(), &newOrder)

		w.WriteHeader(http.StatusAccepted)
		return
	}

	logger.WithError(err).Debug(errAddOrder)

	if errors.Is(err, entity.ErrInvalidOrderNumber) {
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
}

func (h *httpHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(auth.AuthHeaderName)

	orders, err := h.orderService.GetOrders(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Debug(errReadFromRepo)
		http.Error(w, errReadFromRepo, http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	res := []*orderResponse{}
	for _, order := range orders {
		res = append(
			res,
			&orderResponse{
				ID:         order.ID.String(),
				Status:     string(order.Status),
				Accrual:    order.Accrual,
				UploadedAt: order.UploadedAt.Format(),
			},
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.WithError(err).Debug(errEncodeResponseBody)
		http.Error(w, errEncodeResponseBody, http.StatusInternalServerError)
	}
}
