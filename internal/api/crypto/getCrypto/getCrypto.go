package getCrypto

import (
	"CryptoService/internal/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type ResponseCrypto struct {
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdated  string  `json:"last_updated"`
}

func CryptoGetHandler(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := store.GetLatestCrypto()
		res := make([]ResponseCrypto, len(val))
		i := 0
		for _, v := range val {
			res[i] = ResponseCrypto{
				Symbol:       v.Symbol,
				Name:         v.Name,
				CurrentPrice: v.Price,
				LastUpdated:  v.Time.Format(time.RFC3339),
			}
			i++
		}
		c.JSON(http.StatusOK, gin.H{"cryptos": res})
	}
}

func CryptoSymbolGetHandler(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		val, _ := store.GetLatestCrypto()
		if v, ok := val[symbol]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("symbol %s is not being tracked", symbol).Error()})
		} else {
			var resp ResponseCrypto
			resp = ResponseCrypto{
				Symbol:       v.Symbol,
				Name:         v.Name,
				CurrentPrice: v.Price,
				LastUpdated:  v.Time.Format(time.RFC3339),
			}

			c.JSON(http.StatusOK, resp)
		}
	}
}

type responseGetHistory struct {
	Price float64 `json:"price"`
	Time  string  `json:"timestamp"`
}

func CryptoSymbolGetHistoryHandler(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")

		if val, err := store.GetCrypto(symbol); err != nil {
			if strings.HasPrefix(err.Error(), fmt.Sprintf("symbol %s is not being tracked", symbol)) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			var resp []responseGetHistory
			for _, v := range val {
				resp = append(resp, responseGetHistory{
					Price: v.Price,
					Time:  v.Time.Format(time.RFC3339),
				})
			}
			c.JSON(http.StatusOK, gin.H{"symbol": symbol, "history": resp})
		}
	}
}

func CryptoSymbolGetStatsHandler(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		if val, err := store.GetCryptoStats(symbol); err != nil {
			if strings.HasPrefix(err.Error(), fmt.Sprintf("symbol %s is not being tracked", symbol)) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			v, _ := store.GetLatestCrypto()
			c.JSON(http.StatusOK, gin.H{"symbol": symbol, "current_price": v[symbol].Price, "stats": val})
		}
	}
}
