package domain

import (
	"context"
	"time"
)

type LoginStatus string

const (
	LoginStatusSuccess       LoginStatus = "success"
	LoginStatusPasswordError LoginStatus = "password_error"
	LoginStatusInvalidToken  LoginStatus = "invalid_token"
	LoginStatusRefreshToken  LoginStatus = "refresh_token"
)

type LoginLog struct {
	UserID    string
	UserIP    string
	UserAgent string
	Status    LoginStatus
}
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"` // RefreshToken 만료 시간
}
type AuthRepository interface {
	InsertLoginLog(*LoginLog) error
	InsertToken(*Token) error
	DeleteToken(*Token) error
}

type AuthUseCase interface {
	Login(ctx context.Context, user *User) (*Token, error)
	Logout(ctx context.Context, token *Token) error
}
