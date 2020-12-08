package main

import (
	"summer-solutions/graphql-test-server/pkg/middleware"
	"summer-solutions/graphql-test-server/services/reviews/graph"
	"summer-solutions/graphql-test-server/services/reviews/graph/generated"

	"github.com/gin-gonic/gin"

	"github.com/summer-solutions/spring"
)

func main() {
	spring.New("accounts").Build().
		RunServer(
			4003,
			generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}),
			func(Router *gin.Engine) {
				Router.Use(middleware.Cors())
			},
		)
}
