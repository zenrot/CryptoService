package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	CoingeckoKey string `yaml:"coingeckoKey" required:"true"`
	StoragePath  string `yaml:"storage_path"`
	HttpConfig   `yaml:"http-config"`
}

type HttpConfig struct {
	JwtKey  string `yaml:"jwt_key" required:"true"`
	Address string `yaml:"address" env-default:"localhost:8080"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		log.Fatal("no config path provided")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}

	return &cfg
}
