package configEnv

import (
	"github.com/zenrot/CryptoService/internal/config"
	"os"
)

func MustLoad() *config.Config {
	return &config.Config{
		CoingeckoKey:   os.Getenv("COINGECKO_KEY"),
		StorageType:    os.Getenv("STORAGE_TYPE"),
		AuthorizerType: os.Getenv("AUTHORIZER_TYPE"),
		HttpConfig: config.HttpConfig{
			JwtKey:  os.Getenv("JWT_KEY"),
			Address: os.Getenv("CRYPTO_SERVICE_ADDRESS"),
		},
		PostgresConfig: config.PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Dbname:   os.Getenv("POSTGRES_DATABASE"),
		},
	}
}
