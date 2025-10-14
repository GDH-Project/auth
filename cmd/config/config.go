package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type EnvConfig struct {
	DbUrl     string `env:"DB_URL,required"`
	JWTSecret string `env:"JWT_SECRET,required"`
}

func GetConfig() *EnvConfig {
	config := &EnvConfig{}

	if err := env.Parse(config); err != nil {
		log.Fatalf("env를 파싱할 수 없습니다. %v", err)
	}

	return config
}
