package main

import (
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/services/reviews/graph"
	"summer-solutions/graphql-test-server/services/reviews/graph/generated"

	"github.com/summer-solutions/spring/services"

	"github.com/summer-solutions/spring"
)

func main() {
	spring.New("reviews").
		RegisterDIService(
			services.OrmRegistry(entity.Init),
			services.OrmEngine(),
		).
		RegisterGinMiddleware(
			middleware.Cors,
		).
		RunServer(4003, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}
