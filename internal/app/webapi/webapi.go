package webapi

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/config"
)

type AccrualGetter interface {
	GetAccrual(context.Context, string) (*entity.Accrual, error)
}

type WebAPI struct {
	accrualAddr string
}

func New(cfg *config.Config) *WebAPI {
	return &WebAPI{
		accrualAddr: cfg.AccrualAddr,
	}
}
