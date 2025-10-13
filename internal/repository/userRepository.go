package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/GDH-Project/auth/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func (r userRepository) CheckCanCreate(ctx context.Context, user domain.User) bool {
	u := &domain.User{}
	q := "SELECT id FROM auth.user_with_role WHERE email = $1 OR name = $2 "

	if err := r.db.QueryRow(ctx, q, user.Email, user.Name).Scan(&u.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true
		}
	} else if u.ID == "" {
		return true
	}

	return false
}

func (r userRepository) Find(ctx context.Context, user domain.User) (*domain.User, error) {
	type target struct {
		key   string
		value string
	}
	t := target{}
	if user.ID != "" {
		t.key = "id"
		t.value = user.ID
	} else if user.Email != "" {
		t.key = "email"
		t.value = user.Email
	} else if user.Name != "" {
		t.key = "name"
		t.value = user.Name
	} else {
		return nil, errors.New("user 검색 조건이 유효하지 않습니다. id, email, name 중 한가지를 선택해야 합니다.")
	}

	u := &domain.User{}
	q := fmt.Sprintf("SELECT id, name, email, password, role_name FROM auth.user_with_role WHERE %v = $1", t.key)
	if err := r.db.QueryRow(ctx, q, t.value).
		Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Password,
			&u.Role,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.User{}, nil
		}
		return nil, err
	}

	return u, nil
}

func (r userRepository) Create(ctx context.Context, user domain.User) error {
	if canCreate := r.CheckCanCreate(ctx, user); !canCreate {
		return errors.New("user already exists")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// user 삽입 후 id 추출
	q := `
		INSERT INTO auth.users(name,email,password)
		VALUES ($1,$2,$3) RETURNING id
	`
	err = tx.QueryRow(ctx, q,
		user.Name,
		user.Email,
		user.Password,
	).Scan(&user.ID)
	if err != nil {
		return err
	}

	// user id 와 role id 바인딩
	q = `
		INSERT INTO auth.user_roles(user_id, role_id) 
		VALUES ($1,(SELECT id FROM auth.roles WHERE name = $2))
	`
	_, err = tx.Exec(ctx, q,
		user.ID,
		user.Role)
	if err != nil {
		return err
	}

	// 트랜잭션 커밋
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r userRepository) Update(ctx context.Context, user domain.User) error {
	if user.ID == "" {
		return errors.New("user ID is nil")
	}

	q := "UPDATE auth.users SET name=$1, password=$2 WHERE id=$3 AND deleted_at IS NULL"
	if _, err := r.db.Exec(ctx, q,
		user.Name,
		user.Password,
		user.ID); err != nil {
		return err
	}

	return nil
}

func (r userRepository) Delete(ctx context.Context, id string) error {
	q := "UPDATE auth.users SET deleted_at = NOW() WHERE id = $1;"
	if _, err := r.db.Exec(ctx, q, id); err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}
