package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/apex/log/handlers/text"
	"os"
	"runtime/debug"
	"strings"
	log2 "summer-solutions/graphql-test-server/internal/log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/apex/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func RunService(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	if os.Getenv("DEBUG") == "" {
		log.SetHandler(log2.Default)
		log.SetLevel(log.WarnLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.SetHandler(text.Default)
		log.SetLevel(log.DebugLevel)
	}
	r := gin.New()
	r.Use(ginContextToContextMiddleware())
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.POST("/query", graphqlHandler(server))
	r.GET("/", playgroundHandler())
	panic(r.Run(":" + port))
}

func GinContextFromContext(ctx context.Context) *gin.Context {
	return ctx.Value("GinContextKey").(*gin.Context)
}

func graphqlHandler(server graphql.ExecutableSchema) gin.HandlerFunc {
	h := handler.NewDefaultServer(server)
	h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		request := GinContextFromContext(ctx)
		requestPart := gin.H{"requestMethod": request.Request.Method,
			"requestUrl":   request.Request.URL.String(),
			"requestSize":  123432,
			"responseSize": 12342,
			"serverIp":     "123.33.21.21",
			"latency":      "3.5s",
			"protocol":     "HTTP/1.1",
			"userAgent":    request.Request.UserAgent(),
			"referrer":     request.Request.Referer(),
			"status":       503,
			"remoteIp":     request.ClientIP()}
		l := log.WithField("httpRequest", requestPart)
		l = log.WithField("envs", os.Environ())

		var trace string
		traceHeader := request.Request.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", "test-med-281914", traceParts[0])
		}
		l = log.WithField("logging.googleapis.com/trace", trace)
		var message string
		asErr, is := err.(error)
		if is {
			message = asErr.Error()
		} else {
			message = "panic"
		}
		l.Error(message + "\n" + string(debug.Stack()))
		return errors.New("internal server error")
	})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func ginContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
