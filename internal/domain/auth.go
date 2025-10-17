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
	UserID       string
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"` // RefreshToken 만료 시간
}
type AuthRepository interface {
	InsertLoginLog(ctx context.Context, loginLog *LoginLog) error
	InsertToken(ctx context.Context, token *Token) error
	GetTokenByUserIDAndToken(ctx context.Context, token *Token) (*Token, error)
	DeleteTokenByUserID(ctx context.Context, userID string) error
}

type AuthService interface {
	InsertLoginLog(ctx context.Context, loginLog *LoginLog) error
	GetTokenByUserIDAndToken(ctx context.Context, token *Token) (*Token, error)
	InsertToken(ctx context.Context, token *Token) error
	DeleteTokenByUserID(ctx context.Context, userID string) error
}

type AuthUseCase interface {
	LoginWithPassword(ctx context.Context, email string, password string) (*Token, error)
	Logout(ctx context.Context, accessToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*Token, error)
	ValidateToken(ctx context.Context, accessToken string) (*User, error)
}
