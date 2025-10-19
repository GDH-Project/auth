package service

import (
	"context"

	"github.com/GDH-Project/auth/internal/domain"
	apperror "github.com/GDH-Project/auth/internal/resource/common/app_error"
)

type userService struct {
	r domain.UserRepository
}

func (svc *userService) CheckCanCreate(ctx context.Context, user *domain.User) *apperror.Error {
	return svc.r.CheckCanCreate(ctx, user)
}

func (svc *userService) CreateNewUser(ctx context.Context, user *domain.User) *apperror.Error {
	return svc.r.Create(ctx, user)
}

func (svc *userService) FindUserByEmail(ctx context.Context, email string) (*domain.User, *apperror.Error) {
	return svc.r.Find(ctx, &domain.User{Email: email})
}

func (svc *userService) FindUserByUserId(ctx context.Context, id string) (*domain.User, *apperror.Error) {
	return svc.r.Find(ctx, &domain.User{ID: id})
}

func (svc *userService) DeleteUserById(ctx context.Context, id string) *apperror.Error {
	return svc.r.Delete(ctx, &domain.User{ID: id})
}

func (svc *userService) UpdateUserByUserID(ctx context.Context, id string, user *domain.User) *apperror.Error {
	return svc.r.Update(ctx, &domain.User{ID: id, Password: user.Password, Name: user.Name})
}

func NewUserService(ur domain.UserRepository) domain.UserService {
	return &userService{
		r: ur,
	}
}
