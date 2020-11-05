package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/apex/log/handlers/text"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"os"
	"runtime/debug"
	logLocal "summer-solutions/graphql-test-server/internal/log"
	handlerGoogle "summer-solutions/graphql-test-server/internal/log/handler"
	"time"

	"github.com/apex/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func RunService(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	r := gin.New()
	if os.Getenv("DEBUG") == "" {
		log.SetHandler(handlerGoogle.Default)
		log.SetLevel(log.WarnLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.SetHandler(text.Default)
		log.SetLevel(log.DebugLevel)
		r.Use(gin.Logger())
	}
	r.Use(ginContextToContextMiddleware())
	r.Use(cors.Default())
	r.POST("/query", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server))))
	r.GET("/", playgroundHandler())
	r.GET("/ping", pingHandler())
	panic(r.Run(":" + port))
}

func graphqlHandler(server graphql.ExecutableSchema) gin.HandlerFunc {
	h := handler.New(server)

	h.AddTransport(transport.Options{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New(1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		var message string
		asErr, is := err.(error)
		if is {
			message = asErr.Error()
		} else {
			message = "panic"
		}
		logLocal.FromContext(ctx).Error(message + "\n" + string(debug.Stack()))
		return errors.New("internal server error")
	})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func pingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "PONG")
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
