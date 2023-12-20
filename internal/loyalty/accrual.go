package loyalty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/model"
)

func (s *LoyaltyService) getOrderAccrual(order *model.Order) {
	url := fmt.Sprintf("%s/api/orders/%s", s.cfg.AccrualAddr, order.ID)

	logger.WithFields(logger.Fields{
		"URL": url,
	}).Info("API call to accrual service")

	resp, err := http.Get(url)
	if err != nil {
		logger.WithError(err).Debug("API call to accrual service failed")
		return
	}

	var payload model.Order
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		return
	}

	defer resp.Body.Close()

	switch payload.Status {
	case model.StatusInvalid, model.StatusProcessed:
		err := s.repo.UpdateOrder(
			context.Background(),
			&model.Order{
				ID:      order.ID,
				Status:  payload.Status,
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

		err = s.repo.Enroll(
			context.Background(),
			&model.Enrollment{
				User:        order.User,
				Sum:         payload.Accrual,
				ProcessedAt: time.Now(),
			},
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
