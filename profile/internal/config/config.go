package config

import (
	"flag"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Port         int           `yaml:"port" env:"SERVER_PORT" env-default:"8000"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"30s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"30s"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func Get() (Config, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config")
	flag.Parse()

	if configPath == "" {
		return getFromEnv()
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func getFromEnv() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
