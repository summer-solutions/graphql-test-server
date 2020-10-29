package main

import (
	"summer-solutions/graphql-test-server/internal"
	"summer-solutions/graphql-test-server/services/products/graph"
	"summer-solutions/graphql-test-server/services/products/graph/generated"
)

func main() {
	internal.RunService(4002, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}
