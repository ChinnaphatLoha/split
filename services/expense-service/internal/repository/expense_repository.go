package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/split/services/expense-service/internal/model"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense *model.Expense) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO expenses (id, group_id, payer_id, amount, description, split_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, query,
		expense.ID, expense.GroupID, expense.PayerID, expense.Amount,
		expense.Description, expense.SplitType, expense.CreatedAt, expense.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, split := range expense.Splits {
		splitQuery := `
			INSERT INTO expense_splits (id, expense_id, user_id, amount, percentage)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = tx.ExecContext(ctx, splitQuery,
			split.ID, expense.ID, split.UserID, split.Amount, split.Percentage,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id string) (*model.Expense, error) {
	query := `SELECT id, group_id, payer_id, amount, description, split_type, created_at, updated_at FROM expenses WHERE id = $1`
	expense := &model.Expense{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&expense.ID, &expense.GroupID, &expense.PayerID, &expense.Amount,
		&expense.Description, &expense.SplitType, &expense.CreatedAt, &expense.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	splits, err := r.getSplits(ctx, id)
	if err != nil {
		return nil, err
	}
	expense.Splits = splits

	return expense, nil
}

func (r *ExpenseRepository) ListByGroupID(ctx context.Context, groupID string, limit, offset int) ([]*model.Expense, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM expenses WHERE group_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, groupID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, group_id, payer_id, amount, description, split_type, created_at, updated_at
		FROM expenses WHERE group_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.db.QueryContext(ctx, query, groupID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []*model.Expense
	for rows.Next() {
		e := &model.Expense{}
		if err := rows.Scan(&e.ID, &e.GroupID, &e.PayerID, &e.Amount,
			&e.Description, &e.SplitType, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, err
		}
		splits, err := r.getSplits(ctx, e.ID)
		if err != nil {
			return nil, 0, err
		}
		e.Splits = splits
		expenses = append(expenses, e)
	}

	return expenses, total, rows.Err()
}

func (r *ExpenseRepository) Update(ctx context.Context, expense *model.Expense) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE expenses SET payer_id = $1, amount = $2, description = $3, split_type = $4, updated_at = $5
		WHERE id = $6
	`
	_, err = tx.ExecContext(ctx, query,
		expense.PayerID, expense.Amount, expense.Description,
		expense.SplitType, expense.UpdatedAt, expense.ID,
	)
	if err != nil {
		return err
	}

	// Delete existing splits and recreate
	_, err = tx.ExecContext(ctx, `DELETE FROM expense_splits WHERE expense_id = $1`, expense.ID)
	if err != nil {
		return err
	}

	for _, split := range expense.Splits {
		splitQuery := `
			INSERT INTO expense_splits (id, expense_id, user_id, amount, percentage)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = tx.ExecContext(ctx, splitQuery,
			split.ID, expense.ID, split.UserID, split.Amount, split.Percentage,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ExpenseRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM expense_splits WHERE expense_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM expenses WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ExpenseRepository) GetGroupBalances(ctx context.Context, groupID string) ([]model.UserBalance, []model.BalanceDetail, error) {
	// Calculate net balances: for each expense, the payer is owed money and each split participant owes money
	query := `
		SELECT
			e.payer_id,
			es.user_id,
			es.amount
		FROM expenses e
		JOIN expense_splits es ON e.id = es.expense_id
		WHERE e.group_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Net balance map: positive means user is owed, negative means user owes
	netBalance := make(map[string]float64)
	// Track pairwise debts
	pairDebt := make(map[string]float64) // "from:to" -> amount

	for rows.Next() {
		var payerID, userID string
		var amount float64
		if err := rows.Scan(&payerID, &userID, &amount); err != nil {
			return nil, nil, err
		}

		if payerID != userID {
			netBalance[payerID] += amount
			netBalance[userID] -= amount

			key := fmt.Sprintf("%s:%s", userID, payerID)
			pairDebt[key] += amount
		}
	}

	var balances []model.UserBalance
	for userID, balance := range netBalance {
		balances = append(balances, model.UserBalance{
			UserID:  userID,
			Balance: balance,
		})
	}

	var details []model.BalanceDetail
	for key, amount := range pairDebt {
		var from, to string
		fmt.Sscanf(key, "%s:%s", &from, &to)
		// Split the key properly
		for i, c := range key {
			if c == ':' {
				from = key[:i]
				to = key[i+1:]
				break
			}
		}
		if amount > 0 {
			details = append(details, model.BalanceDetail{
				FromUserID: from,
				ToUserID:   to,
				Amount:     amount,
			})
		}
	}

	return balances, details, rows.Err()
}

func (r *ExpenseRepository) getSplits(ctx context.Context, expenseID string) ([]model.ExpenseSplit, error) {
	query := `SELECT id, expense_id, user_id, amount, percentage FROM expense_splits WHERE expense_id = $1`
	rows, err := r.db.QueryContext(ctx, query, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []model.ExpenseSplit
	for rows.Next() {
		s := model.ExpenseSplit{}
		if err := rows.Scan(&s.ID, &s.ExpenseID, &s.UserID, &s.Amount, &s.Percentage); err != nil {
			return nil, err
		}
		splits = append(splits, s)
	}
	return splits, rows.Err()
}
