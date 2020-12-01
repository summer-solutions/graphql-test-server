package main

import (
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/services/accounts/graph"
	"summer-solutions/graphql-test-server/services/accounts/graph/generated"

	"github.com/summer-solutions/spring/service/registry/request"

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
		global.LogGlobalService(),
		global.OrmConfigGlobalService(entity.Init),
	)
	s.RegisterRequestServices(request.LogRequestService())
}
