package usecase

import (
	"context"
	"user-service/internal/entity"
	"user-service/pkg/errWrap"
)

type OrderRepository interface {
	CreateOrder(context.Context, *entity.Order) error
	OrdersByUser(context.Context, int64) ([]entity.Order, error)
}

type OrderService struct {
	repo OrderRepository
	log  entity.Logger
}

func NewOrderService(repo OrderRepository, log entity.Logger) *OrderService {
	return &OrderService{
		repo: repo,
		log:  log,
	}
}

func (os *OrderService) CreateOrder(ctx context.Context, userID int64, orderNumber string) error {
	if !isValidLuna(orderNumber) {
		os.log.Error(ctx, "invalid order number", "number", orderNumber)
		return errWrap.NewAppError(errWrap.ErrValidation, "invalid order number", nil)
	}

	order := entity.Order{
		UserID: userID,
		Number: orderNumber,
		Status: "NEW",
	}
	if err := os.repo.CreateOrder(ctx, &order); err != nil {
		os.log.Error(ctx, "failed to create order", "number", orderNumber, "error", err)
		return err
	}
	
	return nil
}

func (os *OrderService) OrdersByUser(ctx context.Context, userID int64) ([]entity.Order, error) {
	orders, err := os.repo.OrdersByUser(ctx, userID)
	if err != nil {
		os.log.Error(ctx, "failed to get orders by user", "error", err)
		return nil, err
	}
	os.log.Info(ctx, "orders retrieved", "count", len(orders))
	return orders, nil
}

func isValidLuna(orderNumber string) bool {
	sum := 0
	nDigits := len(orderNumber)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit := int(orderNumber[i] - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (sum % 10) == 0
}
