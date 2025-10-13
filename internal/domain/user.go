package domain

import "context"

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
	Find(ctx context.Context, user User) (*User, error)
	CheckCanCreate(ctx context.Context, user User) bool
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}

type UserUseCase interface {
	FindUser(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id string) error
}
