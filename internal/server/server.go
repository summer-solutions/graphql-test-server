package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/apex/log/handlers/text"
	"os"
	"runtime/debug"
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
		l := log.WithField("@type", "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent")
		request := GinContextFromContext(ctx)
		l = l.WithField("httpRequest", gin.H{"requestMethod": request.Request.Method, "url": request.Request.URL.String(),
			"userAgent": request.Request.UserAgent(), "referrer": request.Request.Referer(), "responseStatusCode": 503, "remoteIp": request.ClientIP()})
		l.Error(string(debug.Stack()))
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
