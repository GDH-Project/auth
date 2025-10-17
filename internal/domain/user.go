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

type UserService interface {
	CreateNewUser(ctx context.Context, user *User) *apperror.Error
	GetUserByEmail(ctx context.Context, email string) (*User, *apperror.Error)
	GetUserByUserId(ctx context.Context, id string) (*User, *apperror.Error)
	DeleteUserById(ctx context.Context, id string) *apperror.Error
	UpdateUserByUserID(ctx context.Context, id string, user *User) *apperror.Error
}

type UserUseCase interface {
	CreateUser(ctx context.Context, user *User) *apperror.Error
	UpdateUserByUserID(ctx context.Context, user *User) *apperror.Error
	DeleteUserByUserIDAndPassword(ctx context.Context, id string, user *User) *apperror.Error
	GetUserByEmail(ctx context.Context, email string) (*User, *apperror.Error)
	GetUserByUserID(ctx context.Context, id string) (*User, *apperror.Error)
}
