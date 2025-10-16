package domain

import (
	"context"

	apperror "github.com/GDH-Project/auth/internal/resource/common/app_error"
)

type Role string

const (
	RoleUser   Role = "user"
	RoleDevice Role = "device"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID       string
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string
	Role     Role
}

type UserRepository interface {
	Find(ctx context.Context, user *User) (*User, *apperror.Error)
	CheckCanCreate(ctx context.Context, user *User) *apperror.Error
	Create(ctx context.Context, user *User) *apperror.Error
	Update(ctx context.Context, user *User) *apperror.Error
	Delete(ctx context.Context, user *User) *apperror.Error
}

type UserUseCase interface {
	FindUser(ctx context.Context, user *User) (*User, *apperror.Error)
	CreateUser(ctx context.Context, user *User) *apperror.Error
	UpdateUser(ctx context.Context, user *User) *apperror.Error
	DeleteUser(ctx context.Context, user *User) *apperror.Error
}
