package main

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/initializer"
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/services/products/graph"
	"summer-solutions/graphql-test-server/services/products/graph/generated"
)

func main() {
	server.NewSpring(
		initHandlers,
		middleware.Orm(middleware.SourceWebAPI),
		middleware.Cors,
	).Run(4001, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}

func initHandlers(s *server.Spring) error {
	s.RegisterInitHandler(
		initializer.ConfigHandler,
		initializer.OrmHandler,
	)

	return nil
}
