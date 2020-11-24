package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"summer-solutions/graphql-test-server/internal/gin"

	"github.com/apex/log"
)

const contextKey = "_log"

func FromContext(ctx context.Context) *log.Entry {
	g := gin.FromContext(ctx)
	l, has := g.Get(contextKey)
	if !has {
		var trace string
		traceHeader := g.Request.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GC_PROJECT_ID"), traceParts[0])
		}
		l = log.WithField("logging.googleapis.com/trace", trace)
		g.Set(contextKey, l)
	}
	return l.(*log.Entry)
}
