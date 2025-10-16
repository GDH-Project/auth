package repository

import (
	"context"

	"github.com/GDH-Project/auth/cmd/config"
	"github.com/GDH-Project/auth/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	JWTSecret []byte
	db        *pgxpool.Pool
}

func (r *authRepository) GetToken(ctx context.Context, t *domain.Token) (*domain.Token, error) {
	token := &domain.Token{}

	q := `
			SELECT user_id, refresh_token, expires_at 
			FROM auth.tokens 
			WHERE user_id = $1 
			  AND refresh_token = $2 
			  AND expires_at > NOW();
`
	if err := r.db.QueryRow(ctx, q,
		t.UserID,
		t.RefreshToken).
		Scan(
			&token.UserID,
			&token.RefreshToken,
			&token.ExpiresAt,
		); err != nil {
		return nil, err
	}

	return token, nil
}

func (r *authRepository) InsertLoginLog(ctx context.Context, loginLog *domain.LoginLog) error {
	q := "INSERT INTO auth.login_log(user_id, ip_address, user_agent, status) VALUES ($1,$2,$3,$4);"

	if _, err := r.db.Exec(ctx, q,
		loginLog.UserID,
		loginLog.UserIP,
		loginLog.UserAgent,
		loginLog.Status); err != nil {
		return err
	}

	return nil
}

func (r *authRepository) InsertToken(ctx context.Context, token *domain.Token) error {
	// 기존 토큰 제거
	if err := r.DeleteToken(ctx, token); err != nil {
		return err
	}

	// 전달된 토큰 생성
	q := "INSERT INTO auth.tokens(user_id, refresh_token, expires_at) VALUES ($1,$2,$3);"
	if _, err := r.db.Exec(ctx, q,
		token.UserID,
		token.RefreshToken,
		token.ExpiresAt,
	); err != nil {
		return err
	}
	return nil
}

func (r *authRepository) DeleteToken(ctx context.Context, token *domain.Token) error {
	q := "DELETE FROM auth.tokens WHERE user_id = $1 AND refresh_token = $2;"
	if _, err := r.db.Exec(ctx, q,
		token.UserID,
		token.RefreshToken,
	); err != nil {
		return err
	}
	return nil
}

func NewAuthRepository(cfg *config.EnvConfig, db *pgxpool.Pool) domain.AuthRepository {
	return &authRepository{
		JWTSecret: []byte(cfg.JWTSecret),
		db:        db,
	}
}
