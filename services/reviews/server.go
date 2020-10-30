package main

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/services/reviews/graph"
	"summer-solutions/graphql-test-server/services/reviews/graph/generated"
)

func main() {
	server.RunService(4003, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}
