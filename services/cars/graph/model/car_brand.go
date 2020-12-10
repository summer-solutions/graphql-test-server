//go:generate go run github.com/vektah/dataloaden CarBrandLoader uint64 *summer-solutions/graphql-test-server/services/cars/graph/model.CarBrand

package model

import (
	"context"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring"
	"summer-solutions/graphql-test-server/pkg/entity"
	"time"
)

type CarBrand struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func CarBrandLoaderForContext(ctx context.Context) *CarBrandLoader {
	return spring.GetServiceForRequestRequired(ctx, "car_brand_loader").(*CarBrandLoader)
}

func BuildCarBrandModel(carBrandEntity *entity.CarBrandEntity) *CarBrand {
	return &CarBrand{ID: int64(carBrandEntity.ID), Name: carBrandEntity.Name}
}

func CarBarLoaderService() *spring.ServiceDefinition {
	return &spring.ServiceDefinition{
		Name:   "car_brand_loader",
		Global: false,
		Build: func(ctn di.Container) (interface{}, error) {
			return &CarBrandLoader{
				wait:     2 * time.Millisecond,
				maxBatch: 100,
				fetch: func(keys []uint64) ([]*CarBrand, []error) {
					carBrands := make([]*CarBrand, len(keys))
					errors := make([]error, len(keys))
					ormEngine := ctn.Get("orm_engine_request").(*orm.Engine)
					var carBrandEntities []*entity.CarBrandEntity
					ormEngine.LoadByIDs(keys, &carBrandEntities)
					for i, carBrandEntity := range carBrandEntities {
						carBrands[i] = BuildCarBrandModel(carBrandEntity)
					}
					return carBrands, errors
				},
			}, nil
		},
	}
}
