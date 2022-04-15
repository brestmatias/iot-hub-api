package tracing

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	XRequestIDKey = "X-Request-ID"
)

func Trace(c *gin.Context, method string) (*gin.Context, string, time.Time) {
	return c, method, time.Now()
}

func Un(c *gin.Context, method string, startTime time.Time) {
	elapsed := time.Since(startTime)
	Log("[method:%s][elapsed:%vms] Trace End.", c, method, elapsed.Milliseconds())
}

func RequestId(ginCtx *gin.Context) string {
	if ginCtx == nil {
		return ""
	}
	return ginCtx.GetHeader(XRequestIDKey)
}

func Log(format string, c *gin.Context, v ...interface{}) {
	log.Printf("[request-id:"+RequestId(c)+"]"+format, v...)
}

func VerboseOn(c *gin.Context) bool {
	return c != nil && c.GetHeader("X-Traced") == "true"
}
