package postCrypto

import (
	"CryptoService/internal/api/crypto/getCrypto"
	"CryptoService/internal/priceUpdater"
	"CryptoService/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type requestPostCrypto struct {
	Symbol string `json:"symbol"`
}

func CryptoPostHandler(store storage.Storage, updater priceUpdater.PriceUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requestPostCrypto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := updater.AddCryptoTracking(req.Symbol); err != nil {
			if strings.HasPrefix(err.Error(), "this coin already exists") {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		val, _ := store.GetLatestCrypto()
		v := val[req.Symbol]
		var resp getCrypto.ResponseCrypto

		resp = getCrypto.ResponseCrypto{
			Symbol:       v.Symbol,
			Name:         v.Name,
			CurrentPrice: v.Price,
			LastUpdated:  v.Time.Format(time.RFC3339),
		}

		c.JSON(http.StatusCreated, gin.H{"crypto": resp})
	}
}
