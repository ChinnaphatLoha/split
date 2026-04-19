package repository

import (
	"context"
	"database/sql"

	"github.com/split/services/settlement-service/internal/model"
)

type SettlementRepository struct {
	db *sql.DB
}

func NewSettlementRepository(db *sql.DB) *SettlementRepository {
	return &SettlementRepository{db: db}
}

func (r *SettlementRepository) Create(ctx context.Context, settlement *model.Settlement) error {
	query := `
		INSERT INTO settlements (id, group_id, from_user_id, to_user_id, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		settlement.ID, settlement.GroupID, settlement.FromUserID, settlement.ToUserID,
		settlement.Amount, settlement.Status, settlement.CreatedAt, settlement.UpdatedAt,
	)
	return err
}

func (r *SettlementRepository) GetByID(ctx context.Context, id string) (*model.Settlement, error) {
	query := `SELECT id, group_id, from_user_id, to_user_id, amount, status, created_at, updated_at FROM settlements WHERE id = $1`
	s := &model.Settlement{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.GroupID, &s.FromUserID, &s.ToUserID,
		&s.Amount, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SettlementRepository) UpdateStatus(ctx context.Context, id string, status model.SettlementStatus) error {
	query := `UPDATE settlements SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *SettlementRepository) ListByGroupID(ctx context.Context, groupID string) ([]*model.Settlement, error) {
	query := `
		SELECT id, group_id, from_user_id, to_user_id, amount, status, created_at, updated_at
		FROM settlements WHERE group_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settlements []*model.Settlement
	for rows.Next() {
		s := &model.Settlement{}
		if err := rows.Scan(&s.ID, &s.GroupID, &s.FromUserID, &s.ToUserID,
			&s.Amount, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		settlements = append(settlements, s)
	}
	return settlements, rows.Err()
}
