package model

import "time"

type Group struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Currency    string        `json:"currency"`
	InviteCode  string        `json:"invite_code"`
	OwnerID     string        `json:"owner_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Members     []GroupMember `json:"members"`
}

type GroupMember struct {
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"` // "owner" or "member"
	JoinedAt time.Time `json:"joined_at"`
}
