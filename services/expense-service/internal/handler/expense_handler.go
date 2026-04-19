package handler

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/split/services/expense-service/internal/model"
	"github.com/split/services/expense-service/internal/service"
	pb "github.com/split/services/expense-service/proto"
)

type ExpenseHandler struct {
	pb.UnimplementedExpenseServiceServer
	svc *service.ExpenseService
}

func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{svc: svc}
}

func (h *ExpenseHandler) CreateExpense(ctx context.Context, req *pb.CreateExpenseRequest) (*pb.CreateExpenseResponse, error) {
	if req.GroupId == "" || req.PayerId == "" || req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "group_id, payer_id, and positive amount are required")
	}

	splitType := protoToSplitType(req.SplitType)
	var inputs []service.SplitInput
	for _, s := range req.Splits {
		inputs = append(inputs, service.SplitInput{
			UserID:     s.UserId,
			Amount:     s.Amount,
			Percentage: s.Percentage,
		})
	}

	expense, err := h.svc.CreateExpense(ctx, req.GroupId, req.PayerId, req.Amount, req.Description, splitType, inputs)
	if err != nil {
		if err == service.ErrInvalidAmount {
			return nil, status.Error(codes.InvalidArgument, "invalid amount")
		}
		if err == service.ErrInvalidSplits {
			return nil, status.Error(codes.InvalidArgument, "invalid splits configuration")
		}
		log.Printf("create expense error: %v", err)
		return nil, status.Error(codes.Internal, "failed to create expense")
	}

	return &pb.CreateExpenseResponse{Expense: expenseToProto(expense)}, nil
}

func (h *ExpenseHandler) GetExpense(ctx context.Context, req *pb.GetExpenseRequest) (*pb.GetExpenseResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "expense id is required")
	}

	expense, err := h.svc.GetExpense(ctx, req.Id)
	if err != nil {
		if err == service.ErrExpenseNotFound {
			return nil, status.Error(codes.NotFound, "expense not found")
		}
		log.Printf("get expense error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get expense")
	}

	return &pb.GetExpenseResponse{Expense: expenseToProto(expense)}, nil
}

func (h *ExpenseHandler) ListGroupExpenses(ctx context.Context, req *pb.ListGroupExpensesRequest) (*pb.ListGroupExpensesResponse, error) {
	if req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id is required")
	}

	expenses, total, err := h.svc.ListGroupExpenses(ctx, req.GroupId, int(req.Limit), int(req.Offset))
	if err != nil {
		log.Printf("list group expenses error: %v", err)
		return nil, status.Error(codes.Internal, "failed to list expenses")
	}

	var pbExpenses []*pb.Expense
	for _, e := range expenses {
		pbExpenses = append(pbExpenses, expenseToProto(e))
	}

	return &pb.ListGroupExpensesResponse{
		Expenses:   pbExpenses,
		TotalCount: int32(total),
	}, nil
}

func (h *ExpenseHandler) UpdateExpense(ctx context.Context, req *pb.UpdateExpenseRequest) (*pb.UpdateExpenseResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "expense id is required")
	}

	splitType := protoToSplitType(req.SplitType)
	var inputs []service.SplitInput
	for _, s := range req.Splits {
		inputs = append(inputs, service.SplitInput{
			UserID:     s.UserId,
			Amount:     s.Amount,
			Percentage: s.Percentage,
		})
	}

	expense, err := h.svc.UpdateExpense(ctx, req.Id, req.PayerId, req.Amount, req.Description, splitType, inputs)
	if err != nil {
		if err == service.ErrExpenseNotFound {
			return nil, status.Error(codes.NotFound, "expense not found")
		}
		if err == service.ErrInvalidSplits {
			return nil, status.Error(codes.InvalidArgument, "invalid splits configuration")
		}
		log.Printf("update expense error: %v", err)
		return nil, status.Error(codes.Internal, "failed to update expense")
	}

	return &pb.UpdateExpenseResponse{Expense: expenseToProto(expense)}, nil
}

func (h *ExpenseHandler) DeleteExpense(ctx context.Context, req *pb.DeleteExpenseRequest) (*pb.DeleteExpenseResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "expense id is required")
	}

	err := h.svc.DeleteExpense(ctx, req.Id)
	if err != nil {
		if err == service.ErrExpenseNotFound {
			return nil, status.Error(codes.NotFound, "expense not found")
		}
		log.Printf("delete expense error: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete expense")
	}

	return &pb.DeleteExpenseResponse{Success: true}, nil
}

func (h *ExpenseHandler) GetGroupBalances(ctx context.Context, req *pb.GetGroupBalancesRequest) (*pb.GetGroupBalancesResponse, error) {
	if req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "group_id is required")
	}

	balances, details, err := h.svc.GetGroupBalances(ctx, req.GroupId)
	if err != nil {
		log.Printf("get group balances error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get balances")
	}

	var pbBalances []*pb.UserBalance
	for _, b := range balances {
		pbBalances = append(pbBalances, &pb.UserBalance{
			UserId:  b.UserID,
			Balance: b.Balance,
		})
	}

	var pbDetails []*pb.BalanceDetail
	for _, d := range details {
		pbDetails = append(pbDetails, &pb.BalanceDetail{
			FromUserId: d.FromUserID,
			ToUserId:   d.ToUserID,
			Amount:     d.Amount,
		})
	}

	return &pb.GetGroupBalancesResponse{
		Balances: pbBalances,
		Details:  pbDetails,
	}, nil
}

func protoToSplitType(st pb.SplitType) model.SplitType {
	switch st {
	case pb.SplitType_SPLIT_TYPE_EQUAL:
		return model.SplitTypeEqual
	case pb.SplitType_SPLIT_TYPE_CUSTOM:
		return model.SplitTypeCustom
	case pb.SplitType_SPLIT_TYPE_PERCENTAGE:
		return model.SplitTypePercentage
	default:
		return model.SplitTypeEqual
	}
}

func splitTypeToProto(st model.SplitType) pb.SplitType {
	switch st {
	case model.SplitTypeEqual:
		return pb.SplitType_SPLIT_TYPE_EQUAL
	case model.SplitTypeCustom:
		return pb.SplitType_SPLIT_TYPE_CUSTOM
	case model.SplitTypePercentage:
		return pb.SplitType_SPLIT_TYPE_PERCENTAGE
	default:
		return pb.SplitType_SPLIT_TYPE_UNSPECIFIED
	}
}

func expenseToProto(e *model.Expense) *pb.Expense {
	pbExpense := &pb.Expense{
		Id:          e.ID,
		GroupId:     e.GroupID,
		PayerId:     e.PayerID,
		Amount:      e.Amount,
		Description: e.Description,
		SplitType:   splitTypeToProto(e.SplitType),
		CreatedAt:   timestamppb.New(e.CreatedAt),
		UpdatedAt:   timestamppb.New(e.UpdatedAt),
	}

	for _, s := range e.Splits {
		pbExpense.Splits = append(pbExpense.Splits, &pb.ExpenseSplit{
			Id:         s.ID,
			ExpenseId:  s.ExpenseID,
			UserId:     s.UserID,
			Amount:     s.Amount,
			Percentage: s.Percentage,
		})
	}

	return pbExpense
}
