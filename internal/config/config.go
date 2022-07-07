package config

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort                string        `envconfig:"APP_PORT" default:"8000"`
	DbDsn                  string        `envconfig:"DB_DSN" required:"true"`
	ReadTimeout            int           `envconfig:"READ_TIMEOUT" default:"30" required:"true"`
	WriteTimeout           int           `envconfig:"WRITE_TIMEOUT" default:"30" required:"true"`
	ReadHeaderTimeout      int           `envconfig:"READ_HEADER_TIMEOUT" default:"30" required:"true"`
	SecretKey              string        `envconfig:"SECRET_KEY" required:"true"`
	SecretExpiredIn        time.Duration `envconfig:"SECRET_EXPIRED_IN" default:"10m" required:"true"`
	RefreshSecretKey       string        `envconfig:"REFRESH_SECRET_KEY" required:"true"`
	RefreshSecretExpiredIn time.Duration `envconfig:"REFRESH_SECRET_EXPIRED_IN" default:"72h" required:"true"`
}

func GetConfig() (Config, error) {
	var isLocal bool

	flag.BoolVar(&isLocal, "local", false, "Use env vars from 'local.env' (should be in root)")
	flag.Parse()

	if isLocal {
		log.Println("using local config")
		if err := godotenv.Load("local.env"); err != nil {
			return Config{}, fmt.Errorf("error loading 'local.env' file: %w", err)
		}
	}

	config := Config{}

	err := envconfig.Process("", &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
