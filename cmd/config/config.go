package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type EnvConfig struct {
	DbUrl     string `env:"DB_URL,required"`
	JWTSecret string `env:"JWT_SECRET,required"`
}

func GetConfig() *EnvConfig {
	config := &EnvConfig{}

	if err := env.Parse(config); err != nil {
		zap.S().Fatalw("Failed to parse env", "error", err)
	}

	zap.S().Info("env Load OK")
	return config
}
