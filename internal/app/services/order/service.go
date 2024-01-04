package order

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type orderService struct {
	repo repository.OrderRepository
}

func New(repo repository.OrderRepository) *orderService {
	return &orderService{
		repo: repo,
	}
}
