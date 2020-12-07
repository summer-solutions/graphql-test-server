package main

import (
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/services/accounts/graph"
	"summer-solutions/graphql-test-server/services/accounts/graph/generated"

	"github.com/gin-gonic/gin"

	"github.com/summer-solutions/spring"
)

func main() {
	spring.New("accounts").
		RegisterDIService(
			spring.ServiceDefinitionOrmRegistry(entity.Init),
			spring.ServiceDefinitionOrmEngine(),
			spring.ServiceDefinitionOrmEngineForContext(),
		).Build().
		RunServer(
			4001,
			generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}),
			func(Router *gin.Engine) {
				Router.Use(middleware.Cors())
			},
		)
}
