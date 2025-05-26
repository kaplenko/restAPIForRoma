package storage

import (
	"context"
	"time"
	"user-service/internal/entity"
	"user-service/pkg/errWrap"
)

func (s *Storage) Balance(ctx context.Context, userID int64) (*entity.Balance, error) {
	query := `SELECT current, withdrawn
			  FROM balances
			  WHERE user_id = $1`
	var balance entity.Balance
	if err := s.pool.QueryRow(ctx, query, userID).Scan(&balance.Current, &balance.Withdraw); err != nil {
		s.log.Error(ctx, "failed to get balance", "error", err)
		return nil, errWrap.WrapError(err)
	}
	return &balance, nil
}

func (s *Storage) Withdraw(ctx context.Context, userID, sum int64, orderNumber string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error(ctx, "failed to begin transaction", "error", err)
		return errWrap.WrapError(err)
	}
	defer tx.Rollback(ctx)

	var current int64
	if err = tx.QueryRow(ctx,
		`SELECT current FROM balances WHERE user_id = $1`,
		userID).Scan(&current); err != nil {
		s.log.Error(ctx, "failed to get current balance", "error", err)
		return errWrap.WrapError(err)
	}
	if current < sum {
		s.log.Error(ctx, "insufficient balance for user")
		return errWrap.NewAppError(errWrap.ErrPaymentRequired, "insufficient funds", nil)
	}

	_, err = tx.Exec(ctx,
		`UPDATE balances SET current = current - $1 WHERE user_id = $2`,
		sum, userID)
	if err != nil {
		s.log.Error(ctx, "failed to update balance", "error", err)
		return errWrap.WrapError(err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO withdrawals (user_id, order_number, sum, processed_at)
			 VALUES ($1, $2, $3, $4)`, userID, orderNumber, sum, time.Now())
	if err != nil {
		s.log.Error(ctx, "failed to insert withdrawal", "error", err)
		return errWrap.WrapError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		s.log.Error(ctx, "failed to commit transaction", "error", err)
		return errWrap.WrapError(err)
	}
	return nil
}

func (s *Storage) Withdrawals(ctx context.Context, userID int64) ([]entity.Withdrawal, error) {
	query := `SELECT order_number, sum, processed_at FROM withdrawals WHERE user_id = $1`
	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		s.log.Error(ctx, "failed to get withdrawals", "error", err)
		return nil, errWrap.WrapError(err)
	}
	defer rows.Close()
	var withdraws []entity.Withdrawal
	for rows.Next() {
		var order entity.Withdrawal
		if err = rows.Scan(
			&order.OrderNumber,
			&order.Sum,
			&order.ProcessedAt,
		); err != nil {
			s.log.Error(ctx, "failed to scan withdrawals", "error", err)
			return nil, errWrap.WrapError(err)
		}
		withdraws = append(withdraws, order)
	}
	if len(withdraws) == 0 {
		s.log.Info(ctx, "no withdrawals found")
		return withdraws, errWrap.NewAppError(errWrap.ErrNoResponseData, "no withdrawals found", nil)
	}
	return withdraws, nil
}
