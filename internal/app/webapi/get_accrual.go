package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/logger"
)

type responseAccrual struct {
	ID     string  `json:"order"`
	Status string  `json:"status"`
	Value  float32 `json:"accrual"`
}

func (w *WebAPI) GetAccrual(ctx context.Context, order string) (*entity.Accrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", w.accrualAddr, order)

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
	}, nil
}
