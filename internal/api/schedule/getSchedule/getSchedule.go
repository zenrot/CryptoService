package getSchedule

import (
	"CryptoService/internal/api/schedule"
	"CryptoService/internal/priceUpdater"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func ScheduleGetHandler(updater priceUpdater.PriceUpdater) func(c *gin.Context) {
	return func(c *gin.Context) {

		if val := updater.GetUpdateTime(); val == time.Duration(-1) {
			resp := schedule.Response{
				Enabled:         false,
				IntervalSeconds: "0",
				LastUpdated:     updater.GetLastUpdated().Format(time.RFC3339),
				NextUpdate:      "Never",
			}
			c.JSON(http.StatusOK, resp)
		} else {
			resp := schedule.Response{
				Enabled:         true,
				IntervalSeconds: strconv.Itoa(int(updater.GetUpdateTime() / time.Second)),
				LastUpdated:     updater.GetLastUpdated().Format(time.RFC3339),
				NextUpdate:      updater.GetLastUpdated().Add(updater.GetUpdateTime()).Format(time.RFC3339),
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}
