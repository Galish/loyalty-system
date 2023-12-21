package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Galish/loyalty-system/internal/logger"
	"github.com/Galish/loyalty-system/internal/model"
)

type responseAccrual struct {
	ID     string  `json:"order"`
	Status string  `json:"status"`
	Value  float32 `json:"accrual"`
}

func (s *AccrualService) fetchAccrual(req *request) (*model.Accrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", s.addr, req.order)

	logger.WithFields(logger.Fields{
		"URL": url,
	}).Info("API call to accrual service")

	resp, err := http.Get(url)
	if err != nil {
		logger.WithError(err).Debug("API call to accrual service failed")
		return nil, err
	}

	var res responseAccrual
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		return nil, err
	}

	defer resp.Body.Close()

	logger.WithFields(logger.Fields{
		"accrual": res,
	}).Debug("accrual service response")

	return &model.Accrual{
		Order:  model.OrderNumber(res.ID),
		Status: model.Status(res.Status),
		Value:  res.Value,
		User:   req.user,
	}, nil
}

func (s *AccrualService) applyAccrual(accrual *model.Accrual) error {
	ctx := context.Background()

	err := s.orderRepo.UpdateOrder(
		ctx,
		&model.Order{
			ID:      accrual.Order,
			Status:  accrual.Status,
			Accrual: accrual.Value,
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update order")
		return err
	}

	err = s.balanceRepo.Enroll(
		ctx,
		&model.Enrollment{
			User:        accrual.User,
			Sum:         accrual.Value,
			ProcessedAt: time.Now(),
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update balance")
		return err
	}

	return nil
}
