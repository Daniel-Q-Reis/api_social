package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	PG         `yaml:"postgres"`
	JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Address     string `yaml:"address" env-default:"localhost:8080"`
	Timeout     int    `yaml:"timeout" env-default:"5"`
	IdleTimeout int    `yaml:"idle_timeout" env-default:"60"`
}

type PG struct {
	URL string `env:"PG_URL" env-required:"true"`
}

type JWT struct {
	Secret   string `env:"JWT_SECRET" env-required:"true"`
	TokenTTL int    `yaml:"token_ttl" env-default:"1"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("config error: " + err.Error())
	}

	return &cfg
}