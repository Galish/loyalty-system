package psql

import (
	"context"

	"github.com/Galish/loyalty-system/internal/entity"
)

func (s *psqlStore) UserOrders(ctx context.Context, userID string) ([]*entity.Order, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT uuid, status, accrual, uploaded_at, user_id FROM orders WHERE user_id = $1;",
		userID,
	)
	if err != nil {
		return []*entity.Order{}, err
	}

	defer rows.Close()

	var orders []*entity.Order

	for rows.Next() {
		var order entity.Order

		if err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
			&order.User,
		); err != nil {
			return []*entity.Order{}, err
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return []*entity.Order{}, err
	}

	return orders, nil
}
