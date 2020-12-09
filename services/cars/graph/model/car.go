package model

import (
	"context"
)

type Car struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	BrandID int64
}

func (c *Car) Brand(ctx context.Context) *CarBrand {
	return nil
}
