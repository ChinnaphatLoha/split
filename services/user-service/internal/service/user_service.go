package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/split/services/user-service/internal/model"
	"github.com/split/services/user-service/internal/repository"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Register(ctx context.Context, email, name, password, currency string) (*model.User, string, string, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, "", "", ErrEmailAlreadyExists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	if currency == "" {
		currency = "USD"
	}

	now := time.Now()
	user := &model.User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		AvatarURL:    "",
		Currency:     currency,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, "", "", err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, string, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id, name, avatarURL, currency string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if avatarURL != "" {
		user.AvatarURL = avatarURL
	}
	if currency != "" {
		user.Currency = currency
	}
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUsersByIDs(ctx context.Context, ids []string) ([]*model.User, error) {
	return s.repo.GetByIDs(ctx, ids)
}

func (s *UserService) RefreshToken(ctx context.Context, refreshTokenStr string) (string, string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", "", ErrInvalidToken
	}

	user, err := s.repo.GetByID(ctx, claims.Subject)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *UserService) generateAccessToken(user *model.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) generateRefreshToken(user *model.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
