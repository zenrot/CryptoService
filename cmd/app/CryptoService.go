package main

import (
	"CryptoService/internal/auth/internalAuth"
	"CryptoService/internal/config"
	http_server "CryptoService/internal/http-server"
	"CryptoService/internal/priceUpdater/priceUpdaterMultithreaded"
	"CryptoService/internal/storage/ramstore"
	"flag"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "configPath", "config/config.yaml", "provide path to the config file")
}

func main() {
	cfg := config.MustLoad(configPath)
	storage := ramstore.NewRamStorage()
	pu := priceUpdaterMultithreaded.NewPriceUpdaterMultithreaded(cfg, storage)

	//if err := priceUpdater.AddCryptoTracking("BTC"); err != nil {
	//	fmt.Println(err)
	//}
	//if err := priceUpdater.AddCryptoTracking("BTC"); err != nil {
	//	fmt.Println(err)
	//}
	auth := internalAuth.NewAuthorizer(storage, cfg.JwtKey)

	serv := http_server.NewHttpRouter(cfg, storage, pu, auth)
	serv.Start()

}
