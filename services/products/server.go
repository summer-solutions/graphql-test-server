package main

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/pkg/service/registry/global"
	"summer-solutions/graphql-test-server/pkg/service/registry/request"
	"summer-solutions/graphql-test-server/services/products/graph"
	"summer-solutions/graphql-test-server/services/products/graph/generated"
)

func main() {
	server.NewServer(
		initHandlers,
		middleware.Cors,
	).Run(4001, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}

func initHandlers(s *server.Server, _ *server.Def) {
	s.RegisterGlobalServices(
		global.ConfigGlobalService,
		global.OrmConfigGlobalService,
	)

	s.RegisterRequestServices(
		request.OrmEngineRequestService,
	)
}
