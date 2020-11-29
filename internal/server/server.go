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
	"summer-solutions/graphql-test-server/internal/log/provider"
	"time"
)

const ModeLocal = "local"
const ModeDev = "dev"
const ModeDemo = "demo"
const ModeProd = "prod"
const ModeTest = "test"

var ioCGlobalContainer di.Container

type InitHandler func(s *Server, def *Def)
type GinMiddleware func(engine *gin.Engine) error
type Def struct {
	Name  string
	scope string
	Build func(ctn di.Container) (interface{}, error)
	Close func(obj interface{}) error
}

type Server struct {
	mode            string
	initHandlers    []InitHandler
	requestServices []InitHandler
	middlewares     []GinMiddleware
}

func NewServer(handler InitHandler, middlewares ...GinMiddleware) *Server {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = ModeProd
	}

	s := &Server{mode: mode, middlewares: middlewares}

	s.initializeIoCHandlers(handler)

	return s
}

func (s *Server) Run(defaultPort uint, server graphql.ExecutableSchema) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	r := gin.New()

	if !s.IsInProdMode() {
		logLocal.SetProvider(provider.Google) //TODO from config
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

func (s *Server) RegisterGlobalServices(handlers ...InitHandler) {
	s.initHandlers = append(s.initHandlers, handlers...)
}

func (s *Server) RegisterRequestServices(handlers ...InitHandler) {
	s.requestServices = append(s.requestServices, handlers...)
}

func (s *Server) initializeIoCHandlers(handlerRegister InitHandler) {
	ioCBuilder, _ := di.NewBuilder()

	handlerRegister(s, nil)

	for _, callback := range s.initHandlers {
		def := &Def{}

		callback(s, def)
		if def.Name == "" {
			panic("IoC global service is registered without name")
		}

		if def.Build == nil {
			panic("IoC global service is registered without Build function")
		}
		def.scope = di.App

		err := ioCBuilder.Add(di.Def{
			Name:  def.Name,
			Scope: def.scope,
			Build: def.Build,
			Close: def.Close,
		})
		if err != nil {
			panic(err)
		}
	}

	for _, callback := range s.requestServices {
		def := &Def{}

		callback(s, def)
		if def.Name == "" {
			panic("IoC request service is registered without name")
		}

		if def.Build == nil {
			panic("IoC request service is registered without Build function")
		}

		def.scope = di.Request

		err := ioCBuilder.Add(di.Def{
			Name:  def.Name,
			Scope: def.scope,
			Build: def.Build,
			Close: def.Close,
		})
		if err != nil {
			panic(err)
		}
	}

	ioCGlobalContainer = ioCBuilder.Build()
}

func (s *Server) attachMiddlewares(engine *gin.Engine) {
	for _, middleware := range s.middlewares {
		err := middleware(engine)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) IsInLocalMode() bool {
	return s.mode == ModeLocal
}

func (s *Server) IsInProdMode() bool {
	return s.mode == ModeProd
}

func (s *Server) IsInDevMode() bool {
	return s.mode == ModeDev
}

func (s *Server) IsInDemoMode() bool {
	return s.mode == ModeDemo
}

func (s *Server) IsInTestMode() bool {
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

func GetRequestContainer(ctx context.Context) di.Container {
	c := ctx.Value("GinContextKey").(*gin.Context)

	container, has := c.Get("RequestContainer")
	if has {
		return container.(di.Container)
	}

	ioCRequestContainer, err := ioCGlobalContainer.SubContainer()
	c.Set("RequestContainer", ioCRequestContainer)

	if err != nil {
		panic(err)
	}

	return ioCRequestContainer
}

func GetGlobalContainer() di.Container {
	return ioCGlobalContainer
}
