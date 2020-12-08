package model

import (
	"context"
	"summer-solutions/graphql-test-server/pkg/dic"
	"summer-solutions/graphql-test-server/pkg/entity"
)

type CarBrand struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Car struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	BrandID int64
}

func (c *Car) Brand(ctx context.Context) *CarBrand {
	ormEngine := dic.OrmEngineForContext(ctx)
	carBrandEntity := &entity.CarBrandEntity{}
	ormEngine.LoadByID(uint64(c.BrandID), carBrandEntity)
	return &CarBrand{ID: int64(carBrandEntity.ID), Name: carBrandEntity.Name}
}
