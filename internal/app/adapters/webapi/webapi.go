package webapi

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/config"
)

type AccrualGetter interface {
	GetAccrual(context.Context, string) (*entity.Accrual, error)
}

type webAPI struct {
	accrualAddr string
}

func New(cfg *config.Config) *webAPI {
	return &webAPI{
		accrualAddr: cfg.AccrualAddr,
	}
}
