package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/split/services/user-service/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, avatar_url, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Name, user.PasswordHash, user.AvatarURL, user.Currency,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, email, name, password_hash, avatar_url, currency, created_at, updated_at FROM users WHERE id = $1`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.AvatarURL, &user.Currency, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, name, password_hash, avatar_url, currency, created_at, updated_at FROM users WHERE email = $1`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.AvatarURL, &user.Currency, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET name = $1, avatar_url = $2, currency = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.AvatarURL, user.Currency, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepository) GetByIDs(ctx context.Context, ids []string) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(
		`SELECT id, email, name, password_hash, avatar_url, currency, created_at, updated_at FROM users WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash,
			&user.AvatarURL, &user.Currency, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
