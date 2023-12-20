package psql

import (
	"context"

	"github.com/Galish/loyalty-system/internal/model"
)

func (s *psqlStore) UserOrders(ctx context.Context, userID string) ([]*model.Order, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT uuid, status, accrual, uploaded_at, user_id FROM orders WHERE user_id = $1;",
		userID,
	)
	if err != nil {
		return []*model.Order{}, err
	}

	defer rows.Close()

	var orders []*model.Order

	for rows.Next() {
		var order model.Order

		if err := rows.Scan(
			&order.ID,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
			&order.User,
		); err != nil {
			return []*model.Order{}, err
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return []*model.Order{}, err
	}

	return orders, nil
}
