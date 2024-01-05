package order

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type orderUseCase struct {
	repo repository.OrderRepository
}

func New(repo repository.OrderRepository) *orderUseCase {
	return &orderUseCase{
		repo: repo,
	}
}
