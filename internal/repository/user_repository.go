package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/GDH-Project/auth/internal/domain"
	apperror "github.com/GDH-Project/auth/internal/resource/common/app_error"
	"github.com/GDH-Project/auth/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func (r *userRepository) CheckCanCreate(ctx context.Context, user *domain.User) *apperror.Error {
	u := &domain.User{}
	q := "SELECT id FROM auth.user_with_role WHERE email = $1 OR name = $2 "

	if err := r.db.QueryRow(ctx, q, user.Email, user.Name).Scan(&u.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
	} else if u.ID == "" {
		return nil
	}

	return &apperror.Error{
		Code:        apperror.UserAlreadyUsed,
		UserMessage: "이미 사용중인 email 혹은 name 입니다.",
		DevMessage:  fmt.Sprintf("email: %s 혹은 name: %s는 이미 사용중입니다.", user.Email, user.Name),
		StatusCode:  http.StatusBadRequest,
		Cause:       errors.New("이미 사용중인 email 혹은 name 입니다"),
	}
}

func (r *userRepository) Find(ctx context.Context, user *domain.User) (*domain.User, *apperror.Error) {
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
		return nil, &apperror.Error{
			Code:        apperror.ParameterNotMatch,
			UserMessage: "파라미터가 잘못 입려되었습니다.",
			DevMessage:  fmt.Sprintf("user 검색 조건이 유효하지 않습니다. id, email, name 중 한가지를 선택해야 합니다.(%+v)", user),
			StatusCode:  http.StatusBadRequest,
			Cause:       errors.New("파라미터가 잘못 입력되었습니다"),
		}
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
		return nil, &apperror.Error{
			Code:        apperror.UserNotFound,
			UserMessage: "존재하지 않는 사용자입니다.",
			DevMessage:  fmt.Sprintf("%s = %s 에 해당하는 사용자가 존재하지 않습니다.", t.key, t.value),
			StatusCode:  http.StatusBadRequest,
			Cause:       err,
		}

	}

	return u, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) *apperror.Error {
	if err := r.CheckCanCreate(ctx, user); err != nil {
		return err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return &apperror.Error{
			Code:        apperror.InternalServerError,
			UserMessage: "사용자 생성중 알 수 없는 오류가 발생했습니다.",
			DevMessage:  "트랜잭션을 초기화하는 도중 알 수 없는 오류가 발생했습니다.",
			StatusCode:  http.StatusInternalServerError,
			Cause:       err,
		}
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
		return &apperror.Error{
			Code:        apperror.UserCreateError,
			UserMessage: "사용자 생성중 알 수 없는 오류가 발생했습니다.",
			DevMessage:  fmt.Sprintf("email: %s, name: %s 계정 생성중 오류가 발생했습니다.", user.Email, user.Name),
			StatusCode:  http.StatusInternalServerError,
			Cause:       err,
		}
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
		return &apperror.Error{
			Code:        apperror.UserRoleBindFailed,
			UserMessage: "사용자 생성중 알 수 없는 오류가 발생했습니다.",
			DevMessage:  fmt.Sprintf("userID: %s, roleName: %s 권한 할당중 오류가 발생했습니다.", user.ID, user.Role),
			StatusCode:  http.StatusInternalServerError,
			Cause:       err,
		}
	}

	// 트랜잭션 커밋
	if err = tx.Commit(ctx); err != nil {
		return &apperror.Error{
			Code:        apperror.InternalServerError,
			UserMessage: "사용자 생성중 알 수 없는 오류가 발생했습니다.",
			DevMessage: fmt.Sprintf("사용자: %+v 생성 트랜잭션 커밋중 알 수 없는 오류가 발생했습니다.", &domain.User{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Role:  user.Role,
			}),
			StatusCode: http.StatusInternalServerError,
			Cause:      err,
		}
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) *apperror.Error {
	if user.ID == "" {
		return &apperror.Error{
			Code:        apperror.ParameterNotMatch,
			UserMessage: "사용자 정보를 불러오는데 실패했습니다.",
			DevMessage:  "user.ID 가 존재하지 않습니다.",
			StatusCode:  http.StatusBadRequest,
			Cause:       errors.New("사용자 정보를 불러오는데 실패했습니다"),
		}
	}

	q := "UPDATE auth.users SET name=$1, password=$2 WHERE id=$3 AND deleted_at IS NULL"
	if _, err := r.db.Exec(ctx, q,
		user.Name,
		user.Password,
		user.ID); err != nil {
		return &apperror.Error{
			Code:        apperror.UserUpdateFailed,
			UserMessage: "사용자 정보 업데이트중 알 수 없는 오류가 발생했습니다.",
			DevMessage:  fmt.Sprintf("userID: %s 사용자 업데이트중 알 수 없는 오류가 발생했습니다. (%v)", user.ID, err),
			StatusCode:  http.StatusInternalServerError,
			Cause:       err,
		}
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, user *domain.User) *apperror.Error {
	targetUser, err := r.Find(ctx, &domain.User{ID: user.ID})
	if err != nil {
		return err
	}

	if isMatched, _ := util.CheckPasswordHash(user.Password, targetUser.Password); !isMatched {
		return &apperror.Error{
			Code:        apperror.UserPasswordNotMatch,
			UserMessage: "비밀번호가 일치하지 않습니다.",
			StatusCode:  http.StatusBadRequest,
			Cause:       errors.New("비밀번호가 일치하지 않습니다"),
		}
	}

	q := "UPDATE auth.users SET deleted_at = NOW() WHERE id = $1"
	if _, err := r.db.Exec(ctx, q, user.ID); err != nil {
		return &apperror.Error{
			Code:        apperror.InternalServerError,
			UserMessage: "사용자를 제거하지 못했습니다.",
			DevMessage:  fmt.Sprintf("userID: %s 사용자를 제거중 오류가 발생했습니다.", user.ID),
			StatusCode:  http.StatusInternalServerError,
			Cause:       err,
		}
	}
	return nil
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}
