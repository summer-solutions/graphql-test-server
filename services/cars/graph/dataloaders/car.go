//go:generate go run github.com/vektah/dataloaden CarLoader uint64 *summer-solutions/graphql-test-server/services/cars/graph/model.Car

package dataloaders

import (
	"context"
	"summer-solutions/graphql-test-server/pkg/dic"
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/services/cars/graph/model"
	"time"

	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring"
)

func CarLoaderForContext(ctx context.Context) *CarLoader {
	return spring.GetServiceForRequestRequired(ctx, "car_loader").(*CarLoader)
}

func (l *CarLoader) ConvertEntity(carEntity *entity.CarEntity) *model.Car {
	return &model.Car{ID: int64(carEntity.ID), Name: carEntity.Name}
}

func (l *CarLoader) Cars(ctx context.Context, first int) ([]*model.Car, error) {
	ormEngine := dic.OrmEngineForContext(ctx)
	var carEntities []*entity.CarEntity
	ormEngine.Search(orm.NewWhere("1"), orm.NewPager(1, first), &carEntities)
	total := len(carEntities)
	cars := make([]*model.Car, total)
	if total > 0 {
		leader := CarLoaderForContext(ctx)
		for i, carEntity := range carEntities {
			cars[i] = l.ConvertEntity(carEntity)
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
				fetch: func(keys []uint64) ([]*model.Car, []error) {
					cars := make([]*model.Car, len(keys))
					errors := make([]error, len(keys))
					ormEngine := ctn.Get("orm_engine_request").(*orm.Engine)
					var carEntities []*entity.CarEntity
					ormEngine.LoadByIDs(keys, &carEntities)
					for i, carEntity := range carEntities {
						cars[i] = &model.Car{ID: int64(carEntity.ID), Name: carEntity.Name}
					}
					return cars, errors
				},
			}, nil
		},
	}
}
