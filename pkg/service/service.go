package service

import (
	"context"

	"github.com/summer-solutions/spring/service"

	"github.com/summer-solutions/orm"
)

func OrmConfig() orm.ValidatedRegistry {
	return service.GetGlobalContainer().Get("orm_config").(orm.ValidatedRegistry)
}

func OrmEngineContext(ctx context.Context) *orm.Engine {
	return service.GetRequestContainer(ctx).Get("orm_context").(*orm.Engine)
}
