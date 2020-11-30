package main

import (
	"summer-solutions/graphql-test-server/pkg/middleware"
	globalLocal "summer-solutions/graphql-test-server/pkg/service/registry/global"
	"summer-solutions/graphql-test-server/pkg/service/registry/request"
	"summer-solutions/graphql-test-server/services/products/graph"
	"summer-solutions/graphql-test-server/services/products/graph/generated"

	"github.com/summer-solutions/spring"
	"github.com/summer-solutions/spring/service/registry/global"
)

func main() {
	spring.NewServer(
		initHandlers,
		middleware.Cors,
	).Run(4001, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}

func initHandlers(s *spring.Server, _ *spring.Def) {
	s.RegisterGlobalServices(
		global.LogGlobalService,
		global.ConfigGlobalService,
		globalLocal.OrmConfigGlobalService,
	)

	s.RegisterRequestServices(
		request.OrmEngineRequestService,
	)
}
