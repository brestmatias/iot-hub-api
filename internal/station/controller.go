package station

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	StationService StationService
}

func New(stationService StationService) Controller {
	return Controller{
		StationService: stationService,
	}
}

func (c *Controller) DiscoverStations(ginCtx *gin.Context, ctx context.Context) {
	sta := c.StationService.SeekAndSaveOnlineStations(ginCtx)
	if len(*sta) == 0 {
		ginCtx.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	log.Printf("[method:discover_stations][result:{%+v}] Suscessfull ", sta)
	ginCtx.JSON(http.StatusOK, sta)
}
func (c *Controller) DoHandshake(ginCtx *gin.Context, ctx context.Context) {
	c.StationService.DoHandshake(ginCtx)
	ginCtx.Writer.WriteHeader(http.StatusOK)
}
