package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"os"
	"runtime/debug"
	logLocal "summer-solutions/graphql-test-server/internal/log"
	handlerGoogle "summer-solutions/graphql-test-server/internal/log/handler"
	"time"
)

const ModeLocal = "local"
const ModeDev = "dev"
const ModeDemo = "demo"
const ModeProd = "prod"
const ModeTest = "test"

type InitHandler func(s *Spring) error
type GinMiddleware func(engine *gin.Engine) error

var IoCBuilder *di.Builder
var IoC di.Container

type Spring struct {
	mode         string
	initHandlers []InitHandler
	middlewares  []GinMiddleware
}

func NewSpring(handler InitHandler, middlewares ...GinMiddleware) *Spring {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = ModeProd
	}

	s := &Spring{mode: mode, middlewares: middlewares}

	s.initializeIoCHandlers(handler)
	return s
}

func (s *Spring) Run(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	r := gin.New()

	if !s.IsInProdMode() {
		log.SetHandler(handlerGoogle.Default)
		log.SetLevel(log.WarnLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.SetHandler(text.Default)
		log.SetLevel(log.DebugLevel)
		r.Use(gin.Logger())
	}
	r.Use(ginContextToContextMiddleware())

	s.attachMiddlewares(r)

	r.POST("/query", timeout.New(timeout.WithTimeout(10*time.Second), timeout.WithHandler(graphqlHandler(server))))
	r.GET("/", playgroundHandler())
	panic(r.Run(":" + port))
}

func (s *Spring) RegisterInitHandler(handlers ...InitHandler) {
	s.initHandlers = append(s.initHandlers, handlers...)
}

func (s *Spring) initializeIoCHandlers(handlerRegister func(*Spring) error) {
	IoCBuilder, _ = di.NewBuilder()

	err := handlerRegister(s)
	if err != nil {
		panic(err)
	}

	for _, callback := range s.initHandlers {
		err := callback(s)
		if err != nil {
			panic(err)
		}
	}

	IoC = IoCBuilder.Build()
}

func (s *Spring) attachMiddlewares(engine *gin.Engine) {
	for _, middleware := range s.middlewares {
		err := middleware(engine)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Spring) IsInLocalMode() bool {
	return s.mode == ModeLocal
}

func (s *Spring) IsInProdMode() bool {
	return s.mode == ModeProd
}

func (s *Spring) IsInDevMode() bool {
	return s.mode == ModeDev
}

func (s *Spring) IsInDemoMode() bool {
	return s.mode == ModeDemo
}

func (s *Spring) IsInTestMode() bool {
	return s.mode == ModeTest
}

func graphqlHandler(server graphql.ExecutableSchema) gin.HandlerFunc {
	h := handler.New(server)

	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
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
