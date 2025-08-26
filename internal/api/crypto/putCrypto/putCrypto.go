package putCrypto

import (
	"CryptoService/internal/api/crypto/getCrypto"
	"CryptoService/internal/priceUpdater"
	"CryptoService/internal/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func CryptoPutSymbolRefresh(store storage.Storage, updater priceUpdater.PriceUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")

		if err := updater.RefreshPrice(symbol); err != nil {
			if strings.HasPrefix(err.Error(), fmt.Sprintf("symbol %s is not being tracked", symbol)) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		val, _ := store.GetLatestCrypto()
		var resp getCrypto.ResponseCrypto
		resp = getCrypto.ResponseCrypto{
			Symbol:       val[symbol].Symbol,
			Name:         val[symbol].Name,
			CurrentPrice: val[symbol].Price,
			LastUpdated:  val[symbol].Time.Format(time.RFC3339),
		}
		c.JSON(http.StatusOK, gin.H{"crypto": resp})

	}
}
