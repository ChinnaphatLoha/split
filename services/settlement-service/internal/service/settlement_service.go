package service

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/split/services/settlement-service/internal/model"
	"github.com/split/services/settlement-service/internal/repository"
)

var (
	ErrSettlementNotFound = errors.New("settlement not found")
)

type SettlementService struct {
	repo *repository.SettlementRepository
}

func NewSettlementService(repo *repository.SettlementRepository) *SettlementService {
	return &SettlementService{repo: repo}
}

// ComputeSettlements implements the debt simplification algorithm.
// Given net balances for each user, it computes the minimum number
// of transactions needed to settle all debts.
//
// Algorithm:
// 1. Separate users into creditors (positive balance) and debtors (negative balance)
// 2. Sort both lists by absolute amount (descending)
// 3. Greedily match the largest debtor with the largest creditor
// 4. This produces at most N-1 transactions for N participants
func (s *SettlementService) ComputeSettlements(balances []model.BalanceEntry) []model.Transaction {
	type entry struct {
		userID string
		amount float64
	}

	var creditors []entry // positive balance — are owed money
	var debtors []entry   // negative balance — owe money

	for _, b := range balances {
		if b.Amount > 0.01 {
			creditors = append(creditors, entry{userID: b.UserID, amount: b.Amount})
		} else if b.Amount < -0.01 {
			debtors = append(debtors, entry{userID: b.UserID, amount: -b.Amount}) // store as positive
		}
	}

	// Sort descending by amount
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].amount > creditors[j].amount
	})
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].amount > debtors[j].amount
	})

	var transactions []model.Transaction
	ci, di := 0, 0

	for ci < len(creditors) && di < len(debtors) {
		transferAmount := math.Min(creditors[ci].amount, debtors[di].amount)
		transferAmount = math.Round(transferAmount*100) / 100

		if transferAmount > 0.01 {
			transactions = append(transactions, model.Transaction{
				FromUserID: debtors[di].userID,
				ToUserID:   creditors[ci].userID,
				Amount:     transferAmount,
			})
		}

		creditors[ci].amount -= transferAmount
		debtors[di].amount -= transferAmount

		if creditors[ci].amount < 0.01 {
			ci++
		}
		if debtors[di].amount < 0.01 {
			di++
		}
	}

	return transactions
}

func (s *SettlementService) RecordSettlement(ctx context.Context, groupID, fromUserID, toUserID string, amount float64) (*model.Settlement, error) {
	now := time.Now()
	settlement := &model.Settlement{
		ID:         uuid.New().String(),
		GroupID:    groupID,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Status:     model.SettlementStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, settlement); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, settlement.ID)
}

func (s *SettlementService) MarkCompleted(ctx context.Context, id string) (*model.Settlement, error) {
	settlement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSettlementNotFound
		}
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, id, model.SettlementStatusCompleted); err != nil {
		return nil, err
	}

	settlement.Status = model.SettlementStatusCompleted
	return settlement, nil
}

func (s *SettlementService) ListSettlements(ctx context.Context, groupID string) ([]*model.Settlement, error) {
	return s.repo.ListByGroupID(ctx, groupID)
}
