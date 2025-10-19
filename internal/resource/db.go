package resource

import (
	"context"

	c "github.com/GDH-Project/auth/cmd/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func InitDB(cfg *c.EnvConfig) *pgxpool.Pool {
	target := cfg.DbUrl

	config, err := pgxpool.ParseConfig(target)
	if err != nil {
		zap.S().Fatalw("failed to parse DB config",
			"error", err,
		)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		zap.S().Fatalw("failed to connect to DB",
			"error", err,
		)
	}

	if err = pool.Ping(context.Background()); err != nil {
		zap.S().Fatalw("failed to ping DB",
			"error", err,
		)
	}

	zap.S().Info("Database connection pool initialized successfully")

	return pool
}
