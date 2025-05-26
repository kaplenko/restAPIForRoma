package entity

import "time"

type User struct {
	ID       int64  `json:"-"`
	Username string `json:"login"`
	PassHash []byte `json:"password"`
}

type Order struct {
	ID       int64     `json:"-"`
	Number   string    `json:"number"`
	UserID   int64     `json:"-"`
	Status   string    `json:"status"`
	Accrual  *int64    `json:"accrual"`
	UploadAT time.Time `json:"upload_at"`
}

type Balance struct {
	UserID   int64 `json:"-"`
	Current  int64 `json:"current"`
	Withdraw int64 `json:"withdraw"`
}

type Withdrawal struct {
	ID          int64     `json:"-"`
	UserID      int64     `json:"-"`
	OrderNumber string    `json:"order"`
	Sum         int64     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
