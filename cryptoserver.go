package main

import (
	"flag"
	"log"

	"github.com/zenrot/CryptoService/internal/auth/internalAuth"
	"github.com/zenrot/CryptoService/internal/config/configYaml"
	httpServer "github.com/zenrot/CryptoService/internal/http-server"
	"github.com/zenrot/CryptoService/internal/priceUpdater/priceUpdaterMultithreaded"
	"github.com/zenrot/CryptoService/internal/storage"
	"github.com/zenrot/CryptoService/internal/storage/postgresStorage"
	"github.com/zenrot/CryptoService/internal/storage/ramstore"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "configPath", "config/config.yaml", "provide path to the config file")
}

func main() {
	cfg := configYaml.MustLoad(configPath)

	if cfg.StorageType == "postgres" {
		var storeCrypto storage.Crypto
		var storeAuth storage.Auth
		var err error
		storeCrypto, err = postgresStorage.NewCrypto(cfg)
		storeAuth, err = postgresStorage.NewAuth(cfg)
		if err != nil {
			log.Fatal(err)
		}
		pu := priceUpdaterMultithreaded.New(cfg, storeCrypto)

		auth := internalAuth.New(storeAuth, cfg.JwtKey)

		serv := httpServer.New(cfg, storeCrypto, pu, auth)
		serv.Start()
	} else {
		var store storage.AuthCrypto
		store, err := ramstore.NewRamStorage()
		if err != nil {
			log.Fatal(err)
		}
		pu := priceUpdaterMultithreaded.New(cfg, store)

		auth := internalAuth.New(store, cfg.JwtKey)

		serv := httpServer.New(cfg, store, pu, auth)
		serv.Start()
	}

	//if err = storage.RegisterUser("leh", "1234"); err != nil {
	//	log.Fatal(err)
	//}
	//
	//user, err := storage.LoginUser("leh", "1234")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(user)

}
