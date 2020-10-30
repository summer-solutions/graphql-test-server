package main

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/services/products/graph"
	"summer-solutions/graphql-test-server/services/products/graph/generated"
)

func main() {
	server.RunService(4002, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}
