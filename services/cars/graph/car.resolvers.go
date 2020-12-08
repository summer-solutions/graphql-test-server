package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"summer-solutions/graphql-test-server/pkg/dic"
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/services/cars/graph/generated"
	"summer-solutions/graphql-test-server/services/cars/graph/model"

	"github.com/summer-solutions/orm"
)

func (r *queryResolver) Cars(ctx context.Context) ([]*model.Car, error) {
	ormEngine := dic.OrmEngineForContext(ctx)
	var carEntities []*entity.CarEntity
	ormEngine.Search(orm.NewWhere("1"), nil, &carEntities)
	cars := make([]*model.Car, len(carEntities))
	for i, carEntity := range carEntities {
		cars[i] = &model.Car{
			ID:      fmt.Sprintf("%d", carEntity.ID),
			Name:    carEntity.Name,
			BrandID: uint64(carEntity.Brand.ID),
		}
	}
	return cars, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
