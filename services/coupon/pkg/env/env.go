package env

import (
	"arvan-challenge/pkg/env"
)

// declare Config to hold the environment variables with default values
type Config struct {
	Mode     string `env:"MODE" envDefault:"prod"`
	LogLevel string `env:"LOGLEVEL" envDefault:"info"`
	DBHost   string `env:"DB_HOST" envDefault:"localhost"`
	DBPass   string `env:"DB_PASS" envDefault:"postgres"`
	DBUser   string `env:"DB_USER" envDefault:"postgres"`
	DBName   string `env:"DB_Name" envDefault:"coupon"`
	RedisUrl string `env:"REDIS_URI" envDefault:"localhost:6379"`
	NATSUrl  string `env:"NATS_URI" envDefault:"localhost:4222"`
	Port     int    `env:"PORT" envDefault:"5002"`
	DBPort   int    `env:"DB_PORT" envDefault:"5432"`
}

func ParseConfig() *Config {
	return env.ParseConfig[Config]()
}
