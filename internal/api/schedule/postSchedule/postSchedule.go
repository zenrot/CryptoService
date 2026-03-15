package postSchedule

import (
	"github.com/gin-gonic/gin"
	"github.com/zenrot/CryptoService/internal/priceUpdater"
	"time"
)

func SchedulePostRefreshHandler(updater priceUpdater.PriceUpdater) func(c *gin.Context) {
	return func(c *gin.Context) {
		if num, err := updater.RefreshAllPrices(); err != nil {
			c.JSON(500, err.Error())
			return
		} else {
			c.JSON(200, gin.H{"updated_count": num, "timestamp": time.Now().Format(time.RFC3339)})
			return
		}
	}
}
