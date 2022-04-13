package app

import (
	"context"

	"github.com/gin-gonic/gin"
)

func buildRouter(ctx context.Context, router *gin.Engine, app *App) {

	router.POST("/station/discover", func(c *gin.Context) {
		app.StationController.DiscoverStations(c, ctx)
	})

}
