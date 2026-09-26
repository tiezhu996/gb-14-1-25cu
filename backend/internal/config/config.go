package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config 后端全部配置，通过环境变量注入。
type Config struct {
	Port         string `env:"PORT" envDefault:"8080"`
	Env          string `env:"ENV" envDefault:"development"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
	MongoURI     string `env:"MONGO_URI" envDefault:"mongodb://codelearn_user:codelearn_pwd@localhost:27017/codelearn_db?authSource=admin"`
	DBName       string `env:"DB_NAME" envDefault:"codelearn_db"`
	RedisAddr    string `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RedisPass    string `env:"REDIS_PASSWORD" envDefault:""`
	JWTSecret    string `env:"JWT_SECRET" envDefault:"change_me_to_a_long_random_string"`
	JudgeTimeOut int    `env:"JUDGE_TIMEOUT_SECONDS" envDefault:"10"`
	SeedAdmin    string `env:"SEED_ADMIN" envDefault:"true"`
}

// Load 从环境变量加载配置。
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}
