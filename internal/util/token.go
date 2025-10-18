package util

import (
	"net/http"
	"time"

	apperror "github.com/GDH-Project/auth/internal/resource/common/app_error"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(secret []byte, userId string, userRole string) (string, *apperror.Error) {
	expiredTime := time.Now().Add(12 * time.Hour)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		UserID:   userId,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "GDH-auth",
			Subject:   userId,
		},
	}).SignedString(secret)

	if err != nil {
		zap.S().Error("JWT 생성중 알 수 없는 오류가 발생했습니다.", zap.Error(err))
		return "", &apperror.Error{
			Code:        apperror.InternalServerError,
			UserMessage: "토큰을 생성하지 못했습니다.",
			DevMessage:  "JWT 토큰 생성중 알 수 없는 오류가 발생했습니다.",
			StatusCode:  http.StatusInternalServerError,
			Cause:       err,
		}
	}

	return token, nil
}

func VerifyJWTToken(secret []byte, tokenString string) (*JWTClaims, *apperror.Error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil {
		return nil, &apperror.Error{
			Code:        apperror.BadRequest,
			UserMessage: "잘못된 접근입니다.",
			DevMessage:  "토큰을 신뢰할 수 없습니다.",
			StatusCode:  http.StatusNotFound,
			Cause:       nil,
		}
	}

	// TODO 검증
	claims, ok := token.Claims.(*JWTClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, &apperror.Error{
		Code:        apperror.BadRequest,
		UserMessage: "토큰이 유효하지 않습니다.",
		StatusCode:  http.StatusBadRequest,
		Cause:       nil,
	}
}
