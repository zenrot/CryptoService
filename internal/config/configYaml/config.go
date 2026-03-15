package configYaml

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/zenrot/CryptoService/internal/config"
	"log"
	"os"
)

func MustLoad(configPath string) *config.Config {
	if configPath == "" {
		log.Fatal("no config path provided")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var cfg config.Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}

	return &cfg
}
