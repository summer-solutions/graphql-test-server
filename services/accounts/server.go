package main

import (
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/services/accounts/graph"
	"summer-solutions/graphql-test-server/services/accounts/graph/generated"

	"github.com/summer-solutions/spring/services"

	"github.com/summer-solutions/spring"
)

func main() {
	spring.NewServer("accounts").
		RegisterDIService(
			services.Config("../../config/"),
			services.LogGlobal(),
			services.OrmRegistry(entity.Init),
			services.OrmEngine(),
		).
		RegisterGinMiddleware(
			middleware.Cors,
		).
		Run(4001, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}
