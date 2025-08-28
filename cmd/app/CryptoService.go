package main

import (
	"CryptoService/internal/auth/internalAuth"
	"CryptoService/internal/config"
	http_server "CryptoService/internal/http-server"
	"CryptoService/internal/priceUpdater/priceUpdaterMultithreaded"
	"CryptoService/internal/storage"
	"CryptoService/internal/storage/postgresStorage"
	"CryptoService/internal/storage/ramstore"
	"flag"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "configPath", "config/config.yaml", "provide path to the config file")
}

func main() {
	cfg := config.MustLoad(configPath)
	var store storage.Storage
	if cfg.StorageType == "postgres" {
		var err error
		store, err = postgresStorage.NewPostgresStorage(cfg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		store = ramstore.NewRamStorage()
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

	pu := priceUpdaterMultithreaded.NewPriceUpdaterMultithreaded(cfg, store)

	auth := internalAuth.NewAuthorizer(store, cfg.JwtKey)

	serv := http_server.NewHttpRouter(cfg, store, pu, auth)
	serv.Start()

}
