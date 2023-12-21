package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/model"
	"github.com/Galish/loyalty-system/internal/order"
	"github.com/Galish/loyalty-system/internal/repository"
)

const timeLayout = "2006-01-02T15:04:05-07:00"

type orderResponse struct {
	ID         string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

func (h *httpHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Debug("unable to read request body")
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}

	newOrder := model.Order{
		ID:   model.OrderNumber(string(body)),
		User: r.Header.Get(auth.AuthHeaderName),
	}

	err = h.orderService.AddOrder(r.Context(), newOrder)
	if err != nil {
		logger.WithError(err).Debug("unable to add order")

		if errors.Is(err, order.ErrIncorrectOrderNumber) {
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

	go h.accrualService.GetAccrual(&newOrder)

	w.WriteHeader(http.StatusAccepted)
}

func (h *httpHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(auth.AuthHeaderName)

	orders, err := h.orderService.GetOrders(r.Context(), userID)
	if err != nil {
		logger.WithError(err).Debug("unable to read from repository")
		http.Error(w, "unable to read from repository", http.StatusInternalServerError)
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
				UploadedAt: order.UploadedAt.Format(timeLayout),
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
