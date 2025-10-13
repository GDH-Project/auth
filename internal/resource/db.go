package resource

import (
	"context"
	"log"

	c "github.com/GDH-Project/auth/cmd/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg *c.EnvConfig) *pgxpool.Pool {
	target := cfg.DbUrl

	config, err := pgxpool.ParseConfig(target)
	if err != nil {
		log.Fatalf("failed to parse DB config :%v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("failed to connect to DB :%v", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping DB :%v", err)
	}

	log.Println("Database connection pool initialized successfully")

	return pool
}
