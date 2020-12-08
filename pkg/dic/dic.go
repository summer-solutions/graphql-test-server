package dic

import (
	"context"
	"github.com/summer-solutions/orm"
	"github.com/summer-solutions/spring"
)

func OrmEngineForContext(ctx context.Context) *orm.Engine {
	e, _ := spring.DIC().OrmEngineForContext(ctx)
	return e
}
