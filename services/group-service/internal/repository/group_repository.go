package repository

import (
	"context"
	"database/sql"

	"github.com/split/services/group-service/internal/model"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *model.Group) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO groups (id, name, description, currency, invite_code, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, query,
		group.ID, group.Name, group.Description, group.Currency,
		group.InviteCode, group.OwnerID, group.CreatedAt, group.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Add owner as a member
	memberQuery := `INSERT INTO group_members (group_id, user_id, role, joined_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, memberQuery, group.ID, group.OwnerID, "owner", group.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *GroupRepository) GetByID(ctx context.Context, id string) (*model.Group, error) {
	query := `SELECT id, name, description, currency, invite_code, owner_id, created_at, updated_at FROM groups WHERE id = $1`
	group := &model.Group{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID, &group.Name, &group.Description, &group.Currency,
		&group.InviteCode, &group.OwnerID, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	members, err := r.getMembers(ctx, id)
	if err != nil {
		return nil, err
	}
	group.Members = members

	return group, nil
}

func (r *GroupRepository) ListByUserID(ctx context.Context, userID string) ([]*model.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.currency, g.invite_code, g.owner_id, g.created_at, g.updated_at
		FROM groups g
		INNER JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
		ORDER BY g.updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		g := &model.Group{}
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.Currency,
			&g.InviteCode, &g.OwnerID, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		members, err := r.getMembers(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		g.Members = members
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *GroupRepository) Update(ctx context.Context, group *model.Group) error {
	query := `UPDATE groups SET name = $1, description = $2, currency = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, group.Name, group.Description, group.Currency, group.UpdatedAt, group.ID)
	return err
}

func (r *GroupRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM group_members WHERE group_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM groups WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *GroupRepository) AddMember(ctx context.Context, groupID, userID, role string) error {
	query := `INSERT INTO group_members (group_id, user_id, role, joined_at) VALUES ($1, $2, $3, NOW()) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, groupID, userID, role)
	return err
}

func (r *GroupRepository) RemoveMember(ctx context.Context, groupID, userID string) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

func (r *GroupRepository) UpdateInviteCode(ctx context.Context, groupID, code string) error {
	query := `UPDATE groups SET invite_code = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, code, groupID)
	return err
}

func (r *GroupRepository) GetByInviteCode(ctx context.Context, code string) (*model.Group, error) {
	query := `SELECT id, name, description, currency, invite_code, owner_id, created_at, updated_at FROM groups WHERE invite_code = $1`
	group := &model.Group{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&group.ID, &group.Name, &group.Description, &group.Currency,
		&group.InviteCode, &group.OwnerID, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	members, err := r.getMembers(ctx, group.ID)
	if err != nil {
		return nil, err
	}
	group.Members = members
	return group, nil
}

func (r *GroupRepository) getMembers(ctx context.Context, groupID string) ([]model.GroupMember, error) {
	query := `SELECT user_id, role, joined_at FROM group_members WHERE group_id = $1`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.GroupMember
	for rows.Next() {
		m := model.GroupMember{}
		if err := rows.Scan(&m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}
