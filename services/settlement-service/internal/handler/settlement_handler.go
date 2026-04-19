package handler

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/split/services/settlement-service/internal/model"
	"github.com/split/services/settlement-service/internal/service"
	pb "github.com/split/services/settlement-service/proto"
)

type SettlementHandler struct {
	pb.UnimplementedSettlementServiceServer
	svc *service.SettlementService
}

func NewSettlementHandler(svc *service.SettlementService) *SettlementHandler {
	return &SettlementHandler{svc: svc}
}

func (h *SettlementHandler) ComputeSettlements(ctx context.Context, req *pb.ComputeSettlementsRequest) (*pb.ComputeSettlementsResponse, error) {
	if len(req.Balances) == 0 {
		return nil, status.Error(codes.InvalidArgument, "balances are required")
	}

	var balances []model.BalanceEntry
	for _, b := range req.Balances {
		balances = append(balances, model.BalanceEntry{
			UserID: b.UserId,
			Amount: b.Amount,
		})
	}

	transactions := h.svc.ComputeSettlements(balances)

	var pbTransactions []*pb.Transaction
	for _, t := range transactions {
		pbTransactions = append(pbTransactions, &pb.Transaction{
			FromUserId: t.FromUserID,
			ToUserId:   t.ToUserID,
			Amount:     t.Amount,
		})
	}

	return &pb.ComputeSettlementsResponse{
		Transactions: pbTransactions,
	}, nil
}

func (h *SettlementHandler) RecordSettlement(ctx context.Context, req *pb.RecordSettlementRequest) (*pb.RecordSettlementResponse, error) {
	if req.GroupId == "" || req.FromUserId == "" || req.ToUserId == "" || req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "group_id, from_user_id, to_user_id, and positive amount are required")
	}

	settlement, err := h.svc.RecordSettlement(ctx, req.GroupId, req.FromUserId, req.ToUserId, req.Amount)
	if err != nil {
		log.Printf("record settlement error: %v", err)
		return nil, status.Error(codes.Internal, "failed to record settlement")
	}

	return &pb.RecordSettlementResponse{
		Settlement: settlementToProto(settlement),
	}, nil
}

func (h *SettlementHandler) MarkSettlementCompleted(ctx context.Context, req *pb.MarkSettlementCompletedRequest) (*pb.MarkSettlementCompletedResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "settlement id is required")
	}

	settlement, err := h.svc.MarkCompleted(ctx, req.Id)
	if err != nil {
		if err == service.ErrSettlementNotFound {
			return nil, status.Error(codes.NotFound, "settlement not found")
		}
		log.Printf("mark settlement completed error: %v", err)
		return nil, status.Error(codes.Internal, "failed to mark settlement completed")
	}

	return &pb.MarkSettlementCompletedResponse{
		Settlement: settlementToProto(settlement),
	}, nil
}

func (h *SettlementHandler) ListSettlements(ctx context.Context, req *pb.ListSettlementsRequest) (*pb.ListSettlementsResponse, error) {
	if req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id is required")
	}

	settlements, err := h.svc.ListSettlements(ctx, req.GroupId)
	if err != nil {
		log.Printf("list settlements error: %v", err)
		return nil, status.Error(codes.Internal, "failed to list settlements")
	}

	var pbSettlements []*pb.Settlement
	for _, s := range settlements {
		pbSettlements = append(pbSettlements, settlementToProto(s))
	}

	return &pb.ListSettlementsResponse{
		Settlements: pbSettlements,
	}, nil
}

func settlementToProto(s *model.Settlement) *pb.Settlement {
	var pbStatus pb.SettlementStatus
	switch s.Status {
	case model.SettlementStatusPending:
		pbStatus = pb.SettlementStatus_SETTLEMENT_STATUS_PENDING
	case model.SettlementStatusCompleted:
		pbStatus = pb.SettlementStatus_SETTLEMENT_STATUS_COMPLETED
	}

	return &pb.Settlement{
		Id:         s.ID,
		GroupId:    s.GroupID,
		FromUserId: s.FromUserID,
		ToUserId:   s.ToUserID,
		Amount:     s.Amount,
		Status:     pbStatus,
		CreatedAt:  timestamppb.New(s.CreatedAt),
		UpdatedAt:  timestamppb.New(s.UpdatedAt),
	}
}
