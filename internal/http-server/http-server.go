package http_server

import (
	"github.com/gin-gonic/gin"
	"github.com/zenrot/CryptoService/internal/api/auth/postAuth"
	"github.com/zenrot/CryptoService/internal/api/crypto/deleteCrypto"
	"github.com/zenrot/CryptoService/internal/api/crypto/getCrypto"
	"github.com/zenrot/CryptoService/internal/api/crypto/postCrypto"
	"github.com/zenrot/CryptoService/internal/api/crypto/putCrypto"
	"github.com/zenrot/CryptoService/internal/api/middleware/authMiddleware"
	"github.com/zenrot/CryptoService/internal/api/schedule/getSchedule"
	"github.com/zenrot/CryptoService/internal/api/schedule/postSchedule"
	"github.com/zenrot/CryptoService/internal/api/schedule/putSchedule"
	"github.com/zenrot/CryptoService/internal/auth"
	"github.com/zenrot/CryptoService/internal/auth/internalAuth"
	"github.com/zenrot/CryptoService/internal/config"
	"github.com/zenrot/CryptoService/internal/priceUpdater"
	"github.com/zenrot/CryptoService/internal/priceUpdater/priceUpdaterMultithreaded"
	"github.com/zenrot/CryptoService/internal/storage"
	"github.com/zenrot/CryptoService/internal/storage/ramstore"
)

type httpServer struct {
	httpCfg      *config.HttpConfig
	router       *gin.Engine
	store        storage.Crypto
	auth         auth.Authorizer
	priceUpdater priceUpdater.PriceUpdater
}

func NewHttpRouterNoConfig() *httpServer {
	store, err := ramstore.NewRamStorage()
	if err != nil {
		return nil
	}
	jwtKey := "test"
	config := &config.Config{
		CoingeckoKey: "5dg5h35rVUQusTuKFCFwurnF",
		StorageType:  "",
		HttpConfig: config.HttpConfig{
			JwtKey:  jwtKey,
			Address: "localhost:8000",
		},
	}
	return &httpServer{
		httpCfg:      &config.HttpConfig,
		router:       gin.Default(),
		store:        store,
		auth:         internalAuth.New(store, jwtKey),
		priceUpdater: priceUpdaterMultithreaded.New(config, store),
	}
}
func New(cfg *config.Config, store storage.Crypto, updater priceUpdater.PriceUpdater, authorizer auth.Authorizer) *httpServer {
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
