package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server `yaml:"server"`
	Logger `yaml:"logger"`
}

type Server struct {
	Address string `yaml:"address"`
}

type Logger struct {
	Level uint32 `yaml:"level"`
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

	log.Println(*cfg)

	return cfg
}
