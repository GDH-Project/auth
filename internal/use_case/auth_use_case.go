package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/GDH-Project/auth/cmd/config"
	"github.com/GDH-Project/auth/internal/domain"
	apperror "github.com/GDH-Project/auth/internal/resource/common/app_error"
	"github.com/GDH-Project/auth/internal/util"
	"go.uber.org/zap"
)

type authUseCase struct {
	JWTSecret []byte
	authSvc   domain.AuthService
	userSvc   domain.UserService
}

// GenerateRefreshToken
// 32바이트 크기의 refresh token 생성
func (uc *authUseCase) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		zap.S().Debug("refresh 토큰을 생성하는데 오류가 발생했습니다.", zap.Error(err))
		return "", err
	}

	// base64 URL 인코딩을 통해 안전한 문자열로 변경
	token := base64.URLEncoding.EncodeToString(b)
	return token, nil

}

func (uc *authUseCase) generateToken(ctx context.Context, userId string, userRole domain.Role) (*domain.Token, error) {
	var err error
	var appErr *apperror.Error
	refreshTokenExpiredTime := time.Now().Add(time.Hour * 12).UTC() // 리프레쉬 토큰의 유효기간은 12시간

	// 엑세스 토큰 문자열 생성
	accessToken, appErr := util.GenerateJWTToken(uc.JWTSecret, userId, string(userRole))
	if appErr != nil {
		return nil, errors.New("엑세스 토큰 문자열 생성에 실패했습니다")
	}

	// 리프레쉬 토큰 문자열 생성
	refreshToken, err := uc.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 반환할 토큰 구조체
	token := &domain.Token{
		UserID:       userId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenExpiredTime,
	}

	// 	기존 토큰 제거
	_ = uc.authSvc.DeleteTokenByUserID(ctx, userId)

	// 	토큰 DB에 삽입
	err = uc.authSvc.InsertToken(ctx, token)
	if err != nil {
		zap.S().Debug("토큰을 DB에 INSERT 하지 못했습니다.", zap.Error(err))
		return nil, err
	}

	return token, nil

}
func (uc *authUseCase) LoginWithPassword(ctx context.Context, email string, password string, log *domain.LoginLog) (*domain.Token, error) {
	var appErr *apperror.Error
	var err error

	// 유저 조회
	user, appErr := uc.userSvc.FindUserByEmail(ctx, email)
	if appErr != nil {
		zap.S().Debug("if 문 진입",
			"user", user,
		)
		return nil, appErr
	}

	log.UserID = user.ID

	// 비밀번호 확인
	isSuccess, err := util.CheckPasswordHash(password, user.Password)
	if !isSuccess {
		// 로그인 실패 로그 삽입
		log.Status = domain.LoginStatusPasswordInvalid
		_ = uc.authSvc.InsertLoginLog(ctx, log)

		return nil, err
	}

	// 액세스 토큰 및 리프래쉬 토큰 생성
	token, err := uc.generateToken(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// 로그인 성공 로그 삽입
	log.Status = domain.LoginStatusSuccess
	_ = uc.authSvc.InsertLoginLog(ctx, log)

	return token, nil

}

func (uc *authUseCase) Logout(ctx context.Context, accessToken string) error {
	// 토큰 유효성 검증
	token, err := util.VerifyJWTToken(uc.JWTSecret, accessToken)
	if err != nil {
		// TODO : 토큰 오류 로그인 로그 로직 추가
		zap.S().Debug("토큰이 유효하지 않습니다.", zap.Error(err))
		return err
	}

	// 토큰 제거
	_ = uc.authSvc.DeleteTokenByUserID(ctx, token.UserID)

	return nil
}

func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	var err error
	var appErr *apperror.Error
	token, err := uc.authSvc.GetTokenByUserIDOrToken(ctx, &domain.Token{
		RefreshToken: refreshToken,
	})
	if err != nil {
		// TODO : 토큰 오류 로그인 로그 로직 추가
		zap.S().Debug("auc.RefreshToken() 토큰 Get 오류 발생", zap.Error(err))
		return nil, err
	}

	// 토큰의 시간이 만료되었다면
	if !time.Now().Before(token.ExpiresAt) {
		zap.S().Debug("auc.RefreshToken() 토큰이 만료되었습니다.")
		return nil, errors.New("토큰이 만료되었습니다")
	}

	// 유저 정보 조회
	user, appErr := uc.userSvc.FindUserByUserId(ctx, token.UserID)
	if appErr != nil {
		zap.S().Debug("auc.RefreshToken() 토큰의 UserID 사용자를 찾을 수 없습니다.",
			zap.String("userID", token.UserID),
		)
		return nil, errors.New(appErr.Error())
	}

	newToken, err := uc.generateToken(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return newToken, nil

}

func (uc *authUseCase) ValidateToken(ctx context.Context, accessToken string) (*domain.User, error) {
	var appErr *apperror.Error
	token, appErr := util.VerifyJWTToken(uc.JWTSecret, accessToken)
	if appErr != nil {
		zap.S().Debug("토큰이 유효하지 않습니다.", zap.Error(appErr))
		return nil, errors.New(appErr.Error())
	}

	user, appErr := uc.userSvc.FindUserByUserId(ctx, token.UserID)
	if appErr != nil {
		zap.S().Debug("UserID에 해당하는 사용자가 존재하지 않습니다.", zap.Error(appErr), zap.String("userID", token.UserID))
		return nil, errors.New(appErr.Error())
	}

	return user, nil
}

func NewAuthUseCase(cfg *config.EnvConfig, authService domain.AuthService, userService domain.UserService) domain.AuthUseCase {
	return &authUseCase{
		JWTSecret: []byte(cfg.JWTSecret),
		authSvc:   authService,
		userSvc:   userService,
	}
}
