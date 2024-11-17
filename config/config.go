package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	SMTP       `yaml:"smtp"`
}

type SMTP struct {
	Host     string `yaml:"host" env:"SMTP_HOST"`
	Port     int    `yaml:"port" env-default:"587" env:"SMTP_PORT"`
	User     string `yaml:"user" env:"SMTP_USER"`
	Password string `yaml:"pass" env:"SMTP_PASSWORD"`
	From     string `yaml:"from" env:"SMTP_FROM"`
}

type HTTPServer struct {
	Addr        string        `yaml:"addr" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadConfig() *Config {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("Error to read env variables: %s", err)
	}

	err := cleanenv.ReadConfig("config/local.yaml", &config)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	return &config
}
