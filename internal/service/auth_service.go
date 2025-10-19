package service

import (
	"context"

	"github.com/GDH-Project/auth/internal/domain"
)

type authService struct {
	r domain.AuthRepository
}

func (svc *authService) InsertLoginLog(ctx context.Context, loginLog *domain.LoginLog) error {
	return svc.r.InsertLoginLog(ctx, loginLog)
}

func (svc *authService) GetTokenByUserIDOrToken(ctx context.Context, token *domain.Token) (*domain.Token, error) {
	return svc.r.GetTokenByUserIDOrToken(ctx, token)
}

func (svc *authService) InsertToken(ctx context.Context, token *domain.Token) error {
	return svc.r.InsertToken(ctx, token)
}

func (svc *authService) DeleteTokenByUserID(ctx context.Context, userID string) error {
	return svc.r.DeleteTokenByUserID(ctx, userID)
}

func NewAuthService(authRepo domain.AuthRepository) domain.AuthService {
	return &authService{
		r: authRepo,
	}
}
