package station

import (
	"context"
	"iot-hub-api/internal/repository"
	"iot-hub-api/internal/restclient"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	StationRepository   repository.StationRepository
	StationClient       restclient.StationClient
	HubConfigRepository repository.HubConfigRepository
}

func New(stationRepository repository.StationRepository, hubConfigRepository repository.HubConfigRepository, stationClient restclient.StationClient) Controller {
	return Controller{
		StationRepository:   stationRepository,
		StationClient:       stationClient,
		HubConfigRepository: hubConfigRepository,
	}
}

func (c *Controller) DiscoverStations(ginCtx *gin.Context, ctx context.Context) {
	stationService := NewStationService(c.StationRepository, c.HubConfigRepository, c.StationClient)
	sta := stationService.SeekAndSaveOnlineStations(ginCtx)
	if len(*sta) == 0 {
		ginCtx.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	log.Printf("[method:discover_stations][result:{%+v}] Suscessfull ", sta)
	ginCtx.JSON(http.StatusOK, sta)
}
func (c *Controller) DoHandshake(ginCtx *gin.Context, ctx context.Context) {
	stationService := NewStationService(c.StationRepository, c.HubConfigRepository, c.StationClient)
	stationService.DoHandshake(ginCtx)
	ginCtx.Writer.WriteHeader(http.StatusOK)
}
