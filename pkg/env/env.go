package env

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func ParseConfig[K any]() *K {
	// first we check `APP_MODE` to indicate wether we are in debugging or production mode(default production mode)
	appmode := os.Getenv("APP_MODE")

	if appmode == "" || appmode == "prod" {
		godotenv.Load("prod.env")
	} else if appmode == "debug" {
		godotenv.Load("debug.env")
	}

	// declare cfg that we want from environment variables
	var cfg K
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
