//go:generate go run github.com/vektah/dataloaden CarLoader uint64 *summer-solutions/graphql-test-server/services/cars/graph/model.Car

package model

import (
	"context"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring"
	"summer-solutions/graphql-test-server/pkg/dic"
	"summer-solutions/graphql-test-server/pkg/entity"
	"time"
)

type Car struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	BrandID int64
}

func (c *Car) Brand(ctx context.Context) *CarBrand {
	value, err := CarBrandLoaderForContext(ctx).Load(uint64(c.BrandID))
	// TODO fix!!
	if err != nil {
		panic(err)
	}
	return value
}

func CarLoaderForContext(ctx context.Context) *CarLoader {
	return spring.GetServiceForRequestRequired(ctx, "car_loader").(*CarLoader)
}

func BuildCarModel(carEntity *entity.CarEntity) *Car {
	return &Car{ID: int64(carEntity.ID), Name: carEntity.Name, BrandID: int64(carEntity.Brand.ID)}
}

func (l *CarLoader) Cars(ctx context.Context, first int) ([]*Car, error) {
	ormEngine := dic.OrmEngineForContext(ctx)
	var carEntities []*entity.CarEntity
	ormEngine.Search(orm.NewWhere("1"), orm.NewPager(1, first), &carEntities)
	total := len(carEntities)
	cars := make([]*Car, total)
	if total > 0 {
		leader := CarLoaderForContext(ctx)
		for i, carEntity := range carEntities {
			cars[i] = BuildCarModel(carEntity)
			leader.Prime(uint64(carEntity.ID), cars[i])
		}
	}
	return cars, nil
}

func CarLoaderService() *spring.ServiceDefinition {
	return &spring.ServiceDefinition{
		Name:   "car_loader",
		Global: false,
		Build: func(ctn di.Container) (interface{}, error) {
			return &CarLoader{
				wait:     2 * time.Millisecond,
				maxBatch: 100,
				fetch: func(keys []uint64) ([]*Car, []error) {
					cars := make([]*Car, len(keys))
					errors := make([]error, len(keys))
					ormEngine := ctn.Get("orm_engine_request").(*orm.Engine)
					var carEntities []*entity.CarEntity
					ormEngine.LoadByIDs(keys, &carEntities)
					for i, carEntity := range carEntities {
						cars[i] = BuildCarModel(carEntity)
					}
					return cars, errors
				},
			}, nil
		},
	}
}
