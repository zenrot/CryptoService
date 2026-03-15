package postgresStorage

import (
	"github.com/zenrot/CryptoService/internal/config/configYaml"
	"testing"
)

func TestNewPostgresStorage(t *testing.T) {
	_, err := NewPostgresStorage(&configYaml.Config{
		CoingeckoKey: "",
		StorageType:  "",
		HttpConfig:   configYaml.HttpConfig{},
		PostgresConfig: configYaml.PostgresConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "alexey",
			Password: "",
			Dbname:   "postgres",
		},
	})
	if err != nil {
		t.Error("NewPostgresStorage err:", err)
	}
	return
}
