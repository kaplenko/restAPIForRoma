package usecase

import (
	"context"
	"errors"
	"user-service/internal/entity"
)

type BalanceRepository interface {
	Balance(context.Context, int64) (*entity.Balance, error)
	Withdraw(context.Context, int64, int64, string) error
	Withdrawals(context.Context, int64) ([]entity.Withdrawal, error)
}

type BalanceService struct {
	repo BalanceRepository
	log  entity.Logger
}

func NewBalanceService(repo BalanceRepository, log entity.Logger) *BalanceService {
	return &BalanceService{
		repo,
		log,
	}
}

func (bs *BalanceService) Balance(ctx context.Context, userID int64) (*entity.Balance, error) {
	balance, err := bs.repo.Balance(ctx, userID)
	if err != nil {
		bs.log.Error(ctx, "failed to get balance", "error", err)
		return nil, err
	}
	bs.log.Info(ctx, "got balance", "balance", balance)
	return balance, nil
}

func (bs *BalanceService) Withdraw(ctx context.Context, userID, sum int64, orderNumber string) error {
	if sum <= 0 {
		bs.log.Error(ctx, "invalid sum", "sum", sum)
		return errors.New("sum must be positive")
	}

	if err := bs.repo.Withdraw(ctx, userID, sum, orderNumber); err != nil {
		bs.log.Error(ctx, "failed to withdraw", "sum", sum, "orderNumber", orderNumber, "error", err)
		return err
	}
	bs.log.Info(ctx, "withdrawn", "orderNumber", orderNumber, "sum", sum)
	return nil
}

func (bs *BalanceService) Withdrawals(ctx context.Context, userID int64) ([]entity.Withdrawal, error) {
	withdrawals, err := bs.repo.Withdrawals(ctx, userID)
	if err != nil {
		bs.log.Error(ctx, "failed to get withdrawals", "error", err)
		return nil, err
	}
	bs.log.Info(ctx, "got withdrawals", "count", len(withdrawals))
	return withdrawals, nil
}
