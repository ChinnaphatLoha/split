package model

import "time"

type SettlementStatus string

const (
	SettlementStatusPending   SettlementStatus = "PENDING"
	SettlementStatusCompleted SettlementStatus = "COMPLETED"
)

type Settlement struct {
	ID         string           `json:"id"`
	GroupID    string           `json:"group_id"`
	FromUserID string           `json:"from_user_id"`
	ToUserID   string           `json:"to_user_id"`
	Amount     float64          `json:"amount"`
	Status     SettlementStatus `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type Transaction struct {
	FromUserID string  `json:"from_user_id"`
	ToUserID   string  `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

type BalanceEntry struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"` // positive = owed, negative = owes
}
