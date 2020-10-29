package main

import (
	"summer-solutions/graphql-test-server/internal"
	"summer-solutions/graphql-test-server/services/reviews/graph"
	"summer-solutions/graphql-test-server/services/reviews/graph/generated"
)

func main() {
	internal.RunService(4003, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}
