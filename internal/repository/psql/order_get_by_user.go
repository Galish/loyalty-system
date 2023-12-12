package psql

import (
	"context"

	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *psqlStore) GetUserOrders(ctx context.Context, userID string) ([]*repo.Order, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT uuid, status, accrual, uploaded_at, user_id FROM orders WHERE user_id = $1;",
		userID,
	)
	if err != nil {
		return []*repo.Order{}, err
	}

	defer rows.Close()

	var orders []*repo.Order

	for rows.Next() {
		var order repo.Order

		if err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
			&order.User,
		); err != nil {
			return []*repo.Order{}, err
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return []*repo.Order{}, err
	}

	return orders, nil
}
