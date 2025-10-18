package usecase

import (
	"context"
	"net/http"

	"github.com/GDH-Project/auth/internal/domain"
	apperror "github.com/GDH-Project/auth/internal/resource/common/app_error"
	"github.com/GDH-Project/auth/internal/util"
)

type userUseCase struct {
	userSvc domain.UserService
}

func (uc *userUseCase) CreateUser(ctx context.Context, user *domain.User) *apperror.Error {
	// 생성 가능한 유저인지 확인
	if err := uc.userSvc.CheckCanCreate(ctx, user); err != nil {
		return err
	}

	// 해싱 비밀번호 생성
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return &apperror.Error{
			Code:        apperror.InternalServerError,
			UserMessage: "사용자 생성중 알 수 없는 오류가 발생했습니다.",
			DevMessage:  "유저 패스워드 해싱 중 오류가 발생했습니다.",
			StatusCode:  http.StatusInternalServerError,
			Cause:       err,
		}
	}

	createUserData := &domain.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
		Role:     user.Role,
	}

	// 유저 생성
	if err := uc.userSvc.CreateNewUser(ctx, createUserData); err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) UpdateUserByUserID(ctx context.Context, user *domain.User) *apperror.Error {
	// 유저가 존재하는지 확인
	if _, err := uc.userSvc.FindUserByUserId(ctx, user.ID); err != nil {
		return err
	}

	if user.Password != "" {
		if hashedPassword, err := util.HashPassword(user.Password); err != nil {
			return &apperror.Error{
				Code:        apperror.InternalServerError,
				UserMessage: "비밀번호 업데이트에 실패했습니다.",
				StatusCode:  http.StatusInternalServerError,
				Cause:       err,
			}
		} else {
			user.Password = hashedPassword
		}
	}

	// 입력된 유저 정보 업데이트
	if err := uc.userSvc.UpdateUserByUserID(ctx, user.ID, user); err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) DeleteUserByUserIDAndPassword(ctx context.Context, id string, password string) *apperror.Error {
	u, err := uc.userSvc.FindUserByUserId(ctx, id)
	if err != nil {
		return err
	}

	if isMatched, err := util.CheckPasswordHash(password, u.Password); err != nil || !isMatched {
		return &apperror.Error{
			Code:        apperror.UserPasswordNotMatch,
			UserMessage: "비밀번호가 일치하지 않습니다.",
			StatusCode:  http.StatusBadRequest,
			Cause:       err,
		}
	}

	// 유저 삭제에 실패한 경우
	if err := uc.userSvc.DeleteUserById(ctx, u.ID); err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) FindUserByEmail(ctx context.Context, email string) (*domain.User, *apperror.Error) {
	return uc.userSvc.FindUserByEmail(ctx, email)
}

func (uc *userUseCase) FindUserByUserID(ctx context.Context, id string) (*domain.User, *apperror.Error) {
	return uc.userSvc.FindUserByUserId(ctx, id)
}

func NewUserUseCase(uc domain.UserService) domain.UserUseCase {
	return &userUseCase{
		userSvc: uc,
	}
}
