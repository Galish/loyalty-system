package loyalty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/loyalty-system/internal/logger"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *LoyaltyService) getOrderAccrual(order *Order) {
	url := fmt.Sprintf("%s/api/orders/%s", s.cfg.AccrualAddr, order.ID)

	logger.WithFields(logger.Fields{
		"URL": url,
	}).Info("API call to accrual service")

	resp, err := http.Get(url)
	if err != nil {
		logger.WithError(err).Debug("API call to accrual service failed")
		return
	}

	var payload Order
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	defer resp.Body.Close()

	switch payload.Status {
	case StatusInvalid, StatusProcessed:
		err := s.repo.UpdateOrder(
			context.Background(),
			&repo.Order{
				ID:      order.ID,
				Status:  string(payload.Status),
				Accrual: payload.Accrual,
			},
		)
		if err != nil {
			logger.WithError(err).Debug("unable to update order fields")
			return
		}

		if payload.Accrual == 0 {
			return
		}

		err = s.repo.UpdateBalance(
			context.Background(),
			order.User,
			payload.Accrual,
		)
		if err != nil {
			logger.WithError(err).Debug("unable to update balance")
			return
		}

	default:
		go func() {
			s.orderCh <- order
		}()
	}
}
