package main

import (
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/services/cars/graph"
	"summer-solutions/graphql-test-server/services/cars/graph/generated"
	"summer-solutions/graphql-test-server/services/cars/graph/model"

	"github.com/summer-solutions/spring/scripts"

	"github.com/gin-gonic/gin"

	"github.com/summer-solutions/spring"
)

func main() {
	spring.New("cars").
		RegisterDIService(
			spring.ServiceProviderConfigDirectory("../../config"),
			spring.ServiceDefinitionOrmRegistry(entity.Init),
			spring.ServiceDefinitionOrmEngine(),
			spring.ServiceDefinitionOrmEngineForContext(),
			scripts.ORMAlters(),
			model.CarLoaderService(),
			model.CarBarLoaderService(),
		).Build().
		RunServer(
			4005,
			generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}),
			func(Router *gin.Engine) {
				Router.Use(middleware.Cors())
			},
		)
}
