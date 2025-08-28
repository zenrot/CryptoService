package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	CoingeckoKey   string `yaml:"coingeckoKey" required:"true"`
	StorageType    string `yaml:"storage_type"`
	HttpConfig     `yaml:"http-config"`
	PostgresConfig `yaml:"postgres-storage"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" required:"true""`
	Port     string `yaml:"port" required:"true"`
	User     string `yaml:"user" required:"true"`
	Password string `yaml:"password" required:"true"`
	Dbname   string `yaml:"dbname" required:"true"`
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
