package log

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func FromContext(ctx context.Context) *log.Entry {
	return ctx.Value("Log").(*log.Entry)
}

func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var trace string
		traceHeader := c.Request.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GC_PROJECT_ID"), traceParts[0])
		}
		l := log.WithField("logging.googleapis.com/trace", trace)
		ctx := context.WithValue(c.Request.Context(), "Log", l)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
