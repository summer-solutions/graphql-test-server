package internal

import (
	"cloud.google.com/go/logging"
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"os"
	"runtime/debug"

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
		log.SetHandler(json.Default)
		log.SetLevel(log.WarnLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.SetHandler(text.Default)
		log.SetLevel(log.DebugLevel)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.POST("/query", graphqlHandler(server))
	r.GET("/", playgroundHandler())
	panic(r.Run(":" + port))
}

func graphqlHandler(server graphql.ExecutableSchema) gin.HandlerFunc {
	h := handler.NewDefaultServer(server)
	h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		asErr, is := err.(error)
		l := log.WithField("stack", string(debug.Stack())).WithField("Severity", "Error")
		if is {
			l.WithError(asErr).Error(asErr.Error())
		} else {
			l.Errorf("%v", err)
		}

		client, err := logging.NewClient(context.Background(), "test-med-281914")
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()
		logName := "my-log"
		e := logging.Entry{}
		e.Severity = logging.Error
		client.Logger(logName).Log(e)
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
