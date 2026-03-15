package deleteCrypto

import (
	"github.com/gin-gonic/gin"
	"github.com/zenrot/CryptoService/internal/priceUpdater"
	"github.com/zenrot/CryptoService/internal/storage"
	"net/http"
)

func CryptoDeleteSymbolHandler(store storage.Crypto, updater priceUpdater.PriceUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		if err := updater.DeleteCryptoTracking(symbol); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err := store.DeleteCrypto(symbol); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
