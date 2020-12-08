package entity

import (
	"github.com/summer-solutions/orm"
)

type CarBrandEntity struct {
	orm.ORM
	ID   uint32
	Name string `orm:"required"`
}
