package initialize

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	HTTPHost      string `env:"HTTP_HOST" envDefault:"localhost"`
	HTTPPort      string `env:"HTTP_PORT" envDefault:"3000"`
	PGHost        string `env:"PG_HOST" envDefault:"127.0.0.1"`
	PGPort        string `env:"PG_PORT" envDefault:"5432"`
	PGUser        string `env:"PG_USER" envDefault:"test"`
	PGPassword    string `env:"PG_PASSWORD" envDefault:"test"`
	PGDatabase    string `env:"PG_DB" envDefault:"admin"`
	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASS" envDefault:"admin"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"2"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
