package entity

import (
	"github.com/summer-solutions/orm"
)

type CarEntity struct {
	orm.ORM
	ID    uint32
	Name  string `orm:"required"`
	Brand *CarBrandEntity
}
