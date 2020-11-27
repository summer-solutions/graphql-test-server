package entity

import (
	"github.com/summer-solutions/orm"
)

type UserEntity struct {
	orm.ORM    `orm:"log=log_db_pool;table=users;redisCache"`
	ID         uint64
	Email      string `orm:"unique=Email_FakeDelete:1;required"`
	FakeDelete bool   `orm:"unique=Email_FakeDelete:2;index=Type:2"`

	UserEmailIndex *orm.CachedQuery `queryOne:":Email = ?"`
}
