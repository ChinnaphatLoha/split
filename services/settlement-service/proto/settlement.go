package proto

import (
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type SettlementStatus int32

const (
	SettlementStatus_SETTLEMENT_STATUS_UNSPECIFIED SettlementStatus = 0
	SettlementStatus_SETTLEMENT_STATUS_PENDING     SettlementStatus = 1
	SettlementStatus_SETTLEMENT_STATUS_COMPLETED   SettlementStatus = 2
)

type Settlement struct {
	Id         string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	GroupId    string                 `protobuf:"bytes,2,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	FromUserId string                 `protobuf:"bytes,3,opt,name=from_user_id,json=fromUserId,proto3" json:"from_user_id,omitempty"`
	ToUserId   string                 `protobuf:"bytes,4,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`
	Amount     float64                `protobuf:"fixed64,5,opt,name=amount,proto3" json:"amount,omitempty"`
	Status     SettlementStatus       `protobuf:"varint,6,opt,name=status,proto3,enum=settlement.SettlementStatus" json:"status,omitempty"`
	CreatedAt  *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt  *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

type Transaction struct {
	FromUserId string  `protobuf:"bytes,1,opt,name=from_user_id,json=fromUserId,proto3" json:"from_user_id,omitempty"`
	ToUserId   string  `protobuf:"bytes,2,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`
	Amount     float64 `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

type BalanceEntry struct {
	UserId string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount float64 `protobuf:"fixed64,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

type ComputeSettlementsRequest struct {
	GroupId  string          `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	Balances []*BalanceEntry `protobuf:"bytes,2,rep,name=balances,proto3" json:"balances,omitempty"`
}

type ComputeSettlementsResponse struct {
	Transactions []*Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
}

type RecordSettlementRequest struct {
	GroupId    string  `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	FromUserId string  `protobuf:"bytes,2,opt,name=from_user_id,json=fromUserId,proto3" json:"from_user_id,omitempty"`
	ToUserId   string  `protobuf:"bytes,3,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`
	Amount     float64 `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
}

type RecordSettlementResponse struct {
	Settlement *Settlement `protobuf:"bytes,1,opt,name=settlement,proto3" json:"settlement,omitempty"`
}

type MarkSettlementCompletedRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

type MarkSettlementCompletedResponse struct {
	Settlement *Settlement `protobuf:"bytes,1,opt,name=settlement,proto3" json:"settlement,omitempty"`
}

type ListSettlementsRequest struct {
	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
}

type ListSettlementsResponse struct {
	Settlements []*Settlement `protobuf:"bytes,1,rep,name=settlements,proto3" json:"settlements,omitempty"`
}
