package deleteCrypto

import (
	"CryptoService/internal/priceUpdater"
	"CryptoService/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CryptoDeleteSymbolHandler(store storage.Storage, updater priceUpdater.PriceUpdater) gin.HandlerFunc {
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
