package http_server

import (
	"CryptoService/internal/api/auth/postAuth"
	"CryptoService/internal/api/crypto/deleteCrypto"
	"CryptoService/internal/api/crypto/getCrypto"
	"CryptoService/internal/api/crypto/postCrypto"
	"CryptoService/internal/api/crypto/putCrypto"
	"CryptoService/internal/api/middleware/authMiddleware"
	"CryptoService/internal/api/schedule/getSchedule"
	"CryptoService/internal/api/schedule/postSchedule"
	"CryptoService/internal/api/schedule/putSchedule"
	"CryptoService/internal/auth"
	"CryptoService/internal/auth/internalAuth"
	"CryptoService/internal/config"
	"CryptoService/internal/priceUpdater"
	"CryptoService/internal/priceUpdater/priceUpdaterMultithreaded"
	"CryptoService/internal/storage"
	"CryptoService/internal/storage/ramstore"
	"github.com/gin-gonic/gin"
)

type httpServer struct {
	httpCfg      *config.HttpConfig
	router       *gin.Engine
	store        storage.Storage
	auth         auth.Authorizer
	priceUpdater priceUpdater.PriceUpdater
}

func NewHttpRouterNoConfig() *httpServer {
	store := ramstore.NewRamStorage()
	jwtKey := "test"
	config := &config.Config{
		CoingeckoKey: "5dg5h35rVUQusTuKFCFwurnF",
		StoragePath:  "",
		HttpConfig: config.HttpConfig{
			JwtKey:  jwtKey,
			Address: "localhost:8000",
		},
	}
	return &httpServer{
		httpCfg:      &config.HttpConfig,
		router:       gin.Default(),
		store:        store,
		auth:         internalAuth.NewAuthorizer(store, jwtKey),
		priceUpdater: priceUpdaterMultithreaded.NewPriceUpdaterMultithreaded(config, store),
	}
}
func NewHttpRouter(cfg *config.Config, store storage.Storage, updater priceUpdater.PriceUpdater, authorizer auth.Authorizer) *httpServer {
	return &httpServer{
		httpCfg:      &cfg.HttpConfig,
		router:       gin.Default(),
		store:        store,
		auth:         authorizer,
		priceUpdater: updater,
	}
}

func (hs *httpServer) Start() {

	hs.priceUpdater.Start()

	cryptoHandlers := hs.router.Group("/crypto")
	cryptoHandlers.Use(authMiddleware.AuthMiddleware(hs.auth))
	{
		cryptoHandlers.GET("",
			getCrypto.CryptoGetHandler(hs.store))
		cryptoHandlers.GET("/:symbol",
			getCrypto.CryptoSymbolGetHandler(hs.store))
		cryptoHandlers.GET("/:symbol/history",
			getCrypto.CryptoSymbolGetHistoryHandler(hs.store))
		cryptoHandlers.GET("/:symbol/stats",
			getCrypto.CryptoSymbolGetStatsHandler(hs.store))

		cryptoHandlers.POST("",
			postCrypto.CryptoPostHandler(hs.store, hs.priceUpdater))

		cryptoHandlers.PUT("/:symbol/refresh",
			putCrypto.CryptoPutSymbolRefresh(hs.store, hs.priceUpdater))
		cryptoHandlers.DELETE("/:symbol",
			deleteCrypto.CryptoDeleteSymbolHandler(hs.store, hs.priceUpdater))
	}

	authHandlers := hs.router.Group("/auth")
	{
		authHandlers.POST("login", postAuth.LoginHandler(hs.auth))
		authHandlers.POST("register", postAuth.RegisterHandler(hs.auth))
	}

	scheduleHandlers := hs.router.Group("/schedule")
	scheduleHandlers.Use(authMiddleware.AuthMiddleware(hs.auth))
	{
		scheduleHandlers.GET("", getSchedule.ScheduleGetHandler(hs.priceUpdater))
		scheduleHandlers.PUT("", putSchedule.SchedulePutHandler(hs.priceUpdater))
		scheduleHandlers.POST("trigger", postSchedule.SchedulePostRefreshHandler(hs.priceUpdater))
	}

	hs.router.Run(hs.httpCfg.Address)
}
