package putSchedule

import (
	"CryptoService/internal/api/schedule"
	"CryptoService/internal/priceUpdater"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SchedulePutHandler(updater priceUpdater.PriceUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req schedule.Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Enabled == true {
			if req.IntervalSeconds < 10 || req.IntervalSeconds > 3600 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "interval seconds must be between 10 and 3600"})
				return
			}

			if err := updater.ChangeUpdateTime(time.Duration(req.IntervalSeconds)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"enabled": true, "interval_seconds": req.IntervalSeconds})
		} else {
			if err := updater.StopUpdating(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"enabled": false, "interval_seconds": 0})
		}

	}
}
