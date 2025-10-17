package util

import (
	"net/http"
	"time"

	apperror "github.com/GDH-Project/auth/internal/resource/common/app_error"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(secret []byte, userId string, userRole string) (string, *apperror.Error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userId,
		"user_role": userRole,
		"exp":       time.Now().Add(time.Hour * 2).Unix(),
	}).SignedString(secret)

	if err != nil {
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

func VerifyJWTToken(secret []byte, tokenString string) (*jwt.Token, *apperror.Error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil {
		return nil, &apperror.Error{
			Code:        apperror.BadRequest,
			UserMessage: "잘못된 접근입니다.",
			DevMessage:  "토큰을 신뢰할 수 없습니다.",
			StatusCode:  http.StatusNotFound,
			Cause:       nil,
		}
	}

	return token, nil
}
