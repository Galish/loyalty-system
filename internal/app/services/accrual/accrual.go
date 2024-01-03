package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/logger"
)

type responseAccrual struct {
	ID     string  `json:"order"`
	Status string  `json:"status"`
	Value  float32 `json:"accrual"`
}

func (s *AccrualService) fetchAccrual(ctx context.Context, req *request) (*entity.Accrual, error) {
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

	return &entity.Accrual{
		Order:  res.ID,
		Status: entity.Status(res.Status),
		Value:  res.Value,
		User:   req.user,
	}, nil
}

func (s *AccrualService) applyAccrual(ctx context.Context, accrual *entity.Accrual) error {
	err := s.orderRepo.UpdateOrder(
		ctx,
		&entity.Order{
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
		&entity.Enrollment{
			User:        accrual.User,
			Sum:         accrual.Value,
			ProcessedAt: entity.Time(time.Now()),
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update balance")
		return err
	}

	return nil
}
