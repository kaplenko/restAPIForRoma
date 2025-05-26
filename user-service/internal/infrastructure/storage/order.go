package storage

import (
	"context"
	"user-service/internal/entity"
	"user-service/pkg/errWrap"
)

func (s *Storage) CreateOrder(ctx context.Context, order *entity.Order) error {
	var exists bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS (
SELECT 1 FROM orders WHERE user_id = $1 AND order_number = $2)`,
		order.UserID, order.Number).Scan(&exists)
	if err != nil {
		s.log.Error(ctx, "failed to check if order exists for user", "orderNumber", order.Number, "error", err)
		return errWrap.WrapError(err)
	}
	if exists {
		s.log.Info(ctx, "order already exists by user", "orderNumber", order.Number)
		return errWrap.NewAppError(errWrap.ErrOrderAlreadyUploaded, "order already uploaded", err)
	}

	err = s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM orders WHERE order_number = $1)`,
		order.Number).Scan(&exists)
	if err != nil {
		s.log.Error(ctx, "failed to check if order exists", "orderNumber", order.Number, "error", err)
		return errWrap.WrapError(err)
	}
	if exists {
		s.log.Info(ctx, "order already exists by another user", "orderNumber", order.Number)
		return errWrap.NewAppError(errWrap.ErrOrderAlreadyExists, "order already uploaded by another user", err)
	}

	_, err = s.pool.Exec(ctx, `INSERT INTO orders (user_id, order_number, status) VALUES ($1, $2, $3)`,
		order.UserID, order.Number, order.Status)
	if err != nil {
		s.log.Error(ctx, "failed to insert order", "orderNumber", order.Number, "status", order.Status, "error", err)
		return errWrap.WrapError(err)
	}
	return nil
}

func (s *Storage) OrdersByUser(ctx context.Context, userID int64) ([]entity.Order, error) {
	query := `SELECT order_number, status, accrual, uploaded_at FROM orders 
         	  WHERE user_id = $1
         	  ORDER BY uploaded_at DESC`
	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		s.log.Error(ctx, "failed to query orders by user", "error", err)
		return nil, errWrap.WrapError(err)
	}
	defer rows.Close()

	var orders []entity.Order

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadAT); err != nil {
			s.log.Error(ctx, "failed to scan order by user", "error", err)
			return nil, errWrap.WrapError(err)
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		s.log.Info(ctx, "no orders found by user")
		return orders, errWrap.NewAppError(errWrap.ErrNoResponseData, "no orders found", nil)
	}
	return orders, nil
}
