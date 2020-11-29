package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/service"
	"summer-solutions/graphql-test-server/services/products/graph/generated"
	"summer-solutions/graphql-test-server/services/products/graph/model"
)

func (r *queryResolver) TopProducts(ctx context.Context, first *int) ([]*model.Product, error) {
	_ = server.GetGlobalContainer().Get(service.OrmConfigService)
	_ = server.GetRequestContainer(ctx).Get(service.OrmContextService)

	return hats, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
