package mock

import (
	"context"
	"math/rand"
	"user-service/internal/usecase"
)

type AccrualService struct{}

func NewAccrualService() *AccrualService {
	return &AccrualService{}
}

func (a *AccrualService) RequestCalculation(ctx context.Context, orderNumber string) (*usecase.AccrualResponse, error) {
	statuses := [...]string{
		"REGISTERED",
		"PROCESSING",
		"PROCESSED",
		"INVALID",
	}
	status := statuses[rand.Intn(len(statuses))]

	var accrual *int64
	if status == "PROCESSED" {
		sum := int64(rand.Intn(1000) + 10)
		accrual = &sum
	}

	return &usecase.AccrualResponse{
		OrderNumber: orderNumber,
		Status:      status,
		Accrual:     accrual,
	}, nil
}
