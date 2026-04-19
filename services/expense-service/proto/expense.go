package proto

import (
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type SplitType int32

const (
	SplitType_SPLIT_TYPE_UNSPECIFIED SplitType = 0
	SplitType_SPLIT_TYPE_EQUAL      SplitType = 1
	SplitType_SPLIT_TYPE_CUSTOM     SplitType = 2
	SplitType_SPLIT_TYPE_PERCENTAGE SplitType = 3
)

type ExpenseSplit struct {
	Id         string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ExpenseId  string  `protobuf:"bytes,2,opt,name=expense_id,json=expenseId,proto3" json:"expense_id,omitempty"`
	UserId     string  `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount     float64 `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Percentage float64 `protobuf:"fixed64,5,opt,name=percentage,proto3" json:"percentage,omitempty"`
}

type Expense struct {
	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	GroupId     string                 `protobuf:"bytes,2,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	PayerId     string                 `protobuf:"bytes,3,opt,name=payer_id,json=payerId,proto3" json:"payer_id,omitempty"`
	Amount      float64                `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Description string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	SplitType   SplitType              `protobuf:"varint,6,opt,name=split_type,json=splitType,proto3,enum=expense.SplitType" json:"split_type,omitempty"`
	Splits      []*ExpenseSplit        `protobuf:"bytes,7,rep,name=splits,proto3" json:"splits,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

type ExpenseSplitInput struct {
	UserId     string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount     float64 `protobuf:"fixed64,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Percentage float64 `protobuf:"fixed64,3,opt,name=percentage,proto3" json:"percentage,omitempty"`
}

type CreateExpenseRequest struct {
	GroupId     string               `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	PayerId     string               `protobuf:"bytes,2,opt,name=payer_id,json=payerId,proto3" json:"payer_id,omitempty"`
	Amount      float64              `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Description string               `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	SplitType   SplitType            `protobuf:"varint,5,opt,name=split_type,json=splitType,proto3,enum=expense.SplitType" json:"split_type,omitempty"`
	Splits      []*ExpenseSplitInput `protobuf:"bytes,6,rep,name=splits,proto3" json:"splits,omitempty"`
}

type CreateExpenseResponse struct {
	Expense *Expense `protobuf:"bytes,1,opt,name=expense,proto3" json:"expense,omitempty"`
}

type GetExpenseRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

type GetExpenseResponse struct {
	Expense *Expense `protobuf:"bytes,1,opt,name=expense,proto3" json:"expense,omitempty"`
}

type ListGroupExpensesRequest struct {
	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Limit   int32  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset  int32  `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
}

type ListGroupExpensesResponse struct {
	Expenses   []*Expense `protobuf:"bytes,1,rep,name=expenses,proto3" json:"expenses,omitempty"`
	TotalCount int32      `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
}

type UpdateExpenseRequest struct {
	Id          string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PayerId     string               `protobuf:"bytes,2,opt,name=payer_id,json=payerId,proto3" json:"payer_id,omitempty"`
	Amount      float64              `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Description string               `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	SplitType   SplitType            `protobuf:"varint,5,opt,name=split_type,json=splitType,proto3,enum=expense.SplitType" json:"split_type,omitempty"`
	Splits      []*ExpenseSplitInput `protobuf:"bytes,6,rep,name=splits,proto3" json:"splits,omitempty"`
}

type UpdateExpenseResponse struct {
	Expense *Expense `protobuf:"bytes,1,opt,name=expense,proto3" json:"expense,omitempty"`
}

type DeleteExpenseRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

type DeleteExpenseResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

type GetGroupBalancesRequest struct {
	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
}

type UserBalance struct {
	UserId  string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Balance float64 `protobuf:"fixed64,2,opt,name=balance,proto3" json:"balance,omitempty"`
}

type BalanceDetail struct {
	FromUserId string  `protobuf:"bytes,1,opt,name=from_user_id,json=fromUserId,proto3" json:"from_user_id,omitempty"`
	ToUserId   string  `protobuf:"bytes,2,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`
	Amount     float64 `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

type GetGroupBalancesResponse struct {
	Balances []*UserBalance   `protobuf:"bytes,1,rep,name=balances,proto3" json:"balances,omitempty"`
	Details  []*BalanceDetail `protobuf:"bytes,2,rep,name=details,proto3" json:"details,omitempty"`
}
