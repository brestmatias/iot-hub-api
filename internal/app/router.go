package app

import (
	"context"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	requestid "github.com/sumit-tembe/gin-requestid"
)

func buildRouter(ctx context.Context, router *gin.Engine, app *App) {
	// Middlewares
	{
		//recovery middleware
		router.Use(gin.Recovery())
		//middleware which injects a 'RequestID' into the context and header of each request.
		router.Use(requestid.RequestID(func() string {
			return uuid.NewV4().String()
		},
		))
		//middleware which enhance Gin request logger to include 'RequestID'
		router.Use(gin.LoggerWithConfig(requestid.GetLoggerConfig(nil, nil, nil)))
	}

	router.POST("/station/discover", func(c *gin.Context) {
		app.StationController.DiscoverStations(c, ctx)
	})

	router.POST("/station/handshake", func(c *gin.Context) {
		app.StationController.DoHandshake(c, ctx)
	})

	router.GET("/station/interface", func(c *gin.Context) {
		app.StationController.GetAllInterfaces(c, ctx)
	})

	router.POST("/dispatcher/reloadtasks", func(c *gin.Context) {
		app.DispatcherController.ReloadTasks(c, ctx)
	})
}
