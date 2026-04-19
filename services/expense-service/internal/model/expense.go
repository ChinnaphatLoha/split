package model

import "time"

type SplitType string

const (
	SplitTypeEqual      SplitType = "EQUAL"
	SplitTypeCustom     SplitType = "CUSTOM"
	SplitTypePercentage SplitType = "PERCENTAGE"
)

type Expense struct {
	ID          string         `json:"id"`
	GroupID     string         `json:"group_id"`
	PayerID     string         `json:"payer_id"`
	Amount      float64        `json:"amount"`
	Description string         `json:"description"`
	SplitType   SplitType      `json:"split_type"`
	Splits      []ExpenseSplit `json:"splits"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ExpenseSplit struct {
	ID         string  `json:"id"`
	ExpenseID  string  `json:"expense_id"`
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

type UserBalance struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

type BalanceDetail struct {
	FromUserID string  `json:"from_user_id"`
	ToUserID   string  `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}
