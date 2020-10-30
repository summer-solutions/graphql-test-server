package main

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/services/accounts/graph"
	"summer-solutions/graphql-test-server/services/accounts/graph/generated"
)

func main() {
	server.RunService(4001, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
}
