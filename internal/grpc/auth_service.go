package grpc

import (
	"context"
	"errors"

	"github.com/GDH-Project/auth/internal/domain"
	"github.com/GDH-Project/auth/internal/grpc/authpb"
	"go.uber.org/zap"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	authUseCase domain.AuthUseCase
}

func (s *AuthService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	loginLog := makeLoginLog(ctx)
	token, err := s.authUseCase.LoginWithPassword(ctx, req.GetEmail(), req.GetPassword(), loginLog)
	if err != nil {
		zap.S().Infow("로그인 실패", "email", req.Email, zap.Error(err))
		return nil, errors.New("이메일 혹은 패스워드를 확인해주세요")
	}

	zap.S().Debug("로그인 성공",
		"service", "grpc.auth.Login",
		"email", req.Email,
	)

	return &authpb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	err := s.authUseCase.Logout(ctx, req.GetAccessToken())
	if err != nil {
		zap.S().Infow("토큰이 유효하지 않습니다.", zap.Error(err))
		return nil, errors.New("토큰이 유효하지 않습니다")
	}
	return &authpb.LogoutResponse{}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	token, err := s.authUseCase.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		zap.S().Infow("토큰이 유효하지 않습니다.", zap.Error(err))
		return nil, errors.New("토큰이 유효하지 않습니다")
	}

	return &authpb.RefreshTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *AuthService) Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error) {
	user, err := s.authUseCase.ValidateToken(ctx, req.GetAccessToken())
	if err != nil {
		zap.S().Infow("토큰이 유효하지 않습니다.", zap.Error(err))
		return nil, errors.New("토큰이 유효하지 않습니다")
	}

	var userRole authpb.UserRole
	switch user.Role {
	case domain.RoleUser:
		userRole = authpb.UserRole_BASIC_USER
	case domain.RoleDevice:
		userRole = authpb.UserRole_DATA_USER
	case domain.RoleAdmin:
		userRole = authpb.UserRole_ADMIN
	default:
		userRole = authpb.UserRole_BASIC_USER
	}

	return &authpb.ValidateResponse{
		UserId:   user.ID,
		UserRole: userRole,
	}, nil
}

func NewGrpcAuthService(authUseCase domain.AuthUseCase) authpb.AuthServiceServer {
	return &AuthService{
		authUseCase: authUseCase,
	}
}
