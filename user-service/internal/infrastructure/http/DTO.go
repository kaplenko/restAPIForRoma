package http

import (
	"math"
	"time"
)

type userDTO struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func CentsToRubles(cents int64) float64 {
	return float64(cents) / 100
}

func RublesToCents(ruble float64) int64 {
	return int64(math.Round(ruble * 100))
}

func AccrualToRubles(accrual *int64) float64 {
	if accrual != nil {
		return CentsToRubles(*accrual)
	}
	return 0
}

type OrderRequest struct {
	Number string `json:"order_number"`
}

type OrderResponse struct {
	Number   string    `json:"number"`
	Status   string    `json:"status"`
	Accrual  float64   `json:"accrual"`
	UploadAT time.Time `json:"uploaded_at"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawalResponse struct {
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type BalanceResponse struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdraw"`
}
