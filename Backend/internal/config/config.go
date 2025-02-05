package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server `yaml:"server"`
	DB     `yaml:"db"`
	Logger `yaml:"logger"`
	Secret string `yaml:"secret"`
}

type Server struct {
	Address         string        `yaml:"address" env-default:":8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-default:"postgres"`
	Name     string `yaml:"name" env-default:"house_service"`
	Password string `yaml:"password" env-default:"postgres"`
}

type Logger struct {
	Level string `yaml:"level"`
}

func ParseConfig(path string) *Config {
	var cfg *Config

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return cfg
}
