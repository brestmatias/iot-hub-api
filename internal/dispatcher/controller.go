package dispatcher

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	DispatcherService *DispatcherService
}

func NewController(dispatcherService *DispatcherService) Controller {
	return Controller{
		DispatcherService: dispatcherService,
	}
}

func (c *Controller) ReloadTasks(ginCtx *gin.Context, ctx context.Context) {
	c.DispatcherService.LoadTasks()
	log.Println("[method:ReloadTasks]Suscessfull ")
	ginCtx.Writer.WriteHeader(http.StatusOK)
}