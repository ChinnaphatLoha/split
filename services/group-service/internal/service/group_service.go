package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/split/services/group-service/internal/model"
	"github.com/split/services/group-service/internal/repository"
)

var (
	ErrGroupNotFound   = errors.New("group not found")
	ErrNotOwner        = errors.New("only the owner can perform this action")
	ErrInvalidInvite   = errors.New("invalid invite code")
	ErrAlreadyMember   = errors.New("user is already a member")
	ErrCannotRemoveOwner = errors.New("cannot remove the group owner")
)

type GroupService struct {
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CreateGroup(ctx context.Context, name, description, currency, ownerID string) (*model.Group, error) {
	if currency == "" {
		currency = "USD"
	}

	inviteCode, err := generateInviteCode()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	group := &model.Group{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Currency:    currency,
		InviteCode:  inviteCode,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, group.ID)
}

func (s *GroupService) GetGroup(ctx context.Context, id string) (*model.Group, error) {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	return group, nil
}

func (s *GroupService) ListUserGroups(ctx context.Context, userID string) ([]*model.Group, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *GroupService) UpdateGroup(ctx context.Context, id, name, description, currency string) (*model.Group, error) {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	if name != "" {
		group.Name = name
	}
	if description != "" {
		group.Description = description
	}
	if currency != "" {
		group.Currency = currency
	}
	group.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, group); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *GroupService) DeleteGroup(ctx context.Context, id, userID string) error {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGroupNotFound
		}
		return err
	}
	if group.OwnerID != userID {
		return ErrNotOwner
	}
	return s.repo.Delete(ctx, id)
}

func (s *GroupService) AddMember(ctx context.Context, groupID, userID, role string) (*model.Group, error) {
	if role == "" {
		role = "member"
	}
	if err := s.repo.AddMember(ctx, groupID, userID, role); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, groupID)
}

func (s *GroupService) RemoveMember(ctx context.Context, groupID, userID string) (*model.Group, error) {
	group, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group.OwnerID == userID {
		return nil, ErrCannotRemoveOwner
	}
	if err := s.repo.RemoveMember(ctx, groupID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, groupID)
}

func (s *GroupService) GenerateInviteCode(ctx context.Context, groupID, userID string) (string, error) {
	group, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrGroupNotFound
		}
		return "", err
	}
	if group.OwnerID != userID {
		return "", ErrNotOwner
	}

	code, err := generateInviteCode()
	if err != nil {
		return "", err
	}

	if err := s.repo.UpdateInviteCode(ctx, groupID, code); err != nil {
		return "", err
	}
	return code, nil
}

func (s *GroupService) JoinByInviteCode(ctx context.Context, code, userID string) (*model.Group, error) {
	group, err := s.repo.GetByInviteCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidInvite
		}
		return nil, err
	}

	if err := s.repo.AddMember(ctx, group.ID, userID, "member"); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, group.ID)
}

func generateInviteCode() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
