package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"summer-solutions/graphql-test-server/pkg/dic"
	"summer-solutions/graphql-test-server/services/cars/graph/dataloaders"
	"summer-solutions/graphql-test-server/services/cars/graph/generated"
	"summer-solutions/graphql-test-server/services/cars/graph/model"

	"github.com/summer-solutions/orm"
)

func (r *queryResolver) Cars(ctx context.Context, first *int) ([]*model.Car, error) {
	dic.OrmEngineForContext(ctx).EnableQueryDebug(orm.QueryLoggerSourceDB)
	return dataloaders.CarLoaderForContext(ctx).Cars(ctx, *first)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
