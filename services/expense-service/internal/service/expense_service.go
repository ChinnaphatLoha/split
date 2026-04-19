package service

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/split/services/expense-service/internal/model"
	"github.com/split/services/expense-service/internal/repository"
)

var (
	ErrExpenseNotFound = errors.New("expense not found")
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrInvalidSplits   = errors.New("invalid splits configuration")
)

type ExpenseService struct {
	repo *repository.ExpenseRepository
}

func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

type SplitInput struct {
	UserID     string
	Amount     float64
	Percentage float64
}

func (s *ExpenseService) CreateExpense(ctx context.Context, groupID, payerID string, amount float64, description string, splitType model.SplitType, splitInputs []SplitInput) (*model.Expense, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	splits, err := s.calculateSplits(amount, splitType, splitInputs)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expense := &model.Expense{
		ID:          uuid.New().String(),
		GroupID:     groupID,
		PayerID:     payerID,
		Amount:      amount,
		Description: description,
		SplitType:   splitType,
		Splits:      splits,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, expense); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, expense.ID)
}

func (s *ExpenseService) GetExpense(ctx context.Context, id string) (*model.Expense, error) {
	expense, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrExpenseNotFound
		}
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) ListGroupExpenses(ctx context.Context, groupID string, limit, offset int) ([]*model.Expense, int, error) {
	return s.repo.ListByGroupID(ctx, groupID, limit, offset)
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, id, payerID string, amount float64, description string, splitType model.SplitType, splitInputs []SplitInput) (*model.Expense, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrExpenseNotFound
		}
		return nil, err
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	splits, err := s.calculateSplits(amount, splitType, splitInputs)
	if err != nil {
		return nil, err
	}

	existing.PayerID = payerID
	existing.Amount = amount
	existing.Description = description
	existing.SplitType = splitType
	existing.Splits = splits
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrExpenseNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *ExpenseService) GetGroupBalances(ctx context.Context, groupID string) ([]model.UserBalance, []model.BalanceDetail, error) {
	return s.repo.GetGroupBalances(ctx, groupID)
}

func (s *ExpenseService) calculateSplits(totalAmount float64, splitType model.SplitType, inputs []SplitInput) ([]model.ExpenseSplit, error) {
	if len(inputs) == 0 {
		return nil, ErrInvalidSplits
	}

	var splits []model.ExpenseSplit

	switch splitType {
	case model.SplitTypeEqual:
		perPerson := math.Round(totalAmount/float64(len(inputs))*100) / 100
		remainder := math.Round((totalAmount-perPerson*float64(len(inputs)))*100) / 100

		for i, input := range inputs {
			amount := perPerson
			if i == 0 {
				amount += remainder // Give remainder to first person
			}
			splits = append(splits, model.ExpenseSplit{
				ID:         uuid.New().String(),
				UserID:     input.UserID,
				Amount:     amount,
				Percentage: 100.0 / float64(len(inputs)),
			})
		}

	case model.SplitTypeCustom:
		var total float64
		for _, input := range inputs {
			total += input.Amount
		}
		if math.Abs(total-totalAmount) > 0.01 {
			return nil, ErrInvalidSplits
		}
		for _, input := range inputs {
			splits = append(splits, model.ExpenseSplit{
				ID:         uuid.New().String(),
				UserID:     input.UserID,
				Amount:     input.Amount,
				Percentage: (input.Amount / totalAmount) * 100,
			})
		}

	case model.SplitTypePercentage:
		var totalPct float64
		for _, input := range inputs {
			totalPct += input.Percentage
		}
		if math.Abs(totalPct-100.0) > 0.01 {
			return nil, ErrInvalidSplits
		}
		for _, input := range inputs {
			amount := math.Round(totalAmount*input.Percentage/100*100) / 100
			splits = append(splits, model.ExpenseSplit{
				ID:         uuid.New().String(),
				UserID:     input.UserID,
				Amount:     amount,
				Percentage: input.Percentage,
			})
		}

	default:
		return nil, ErrInvalidSplits
	}

	return splits, nil
}
