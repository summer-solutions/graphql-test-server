package service

import (
	"context"
	"summer-solutions/graphql-test-server/internal/service"

	"github.com/summer-solutions/orm"
)

func OrmConfig() orm.ValidatedRegistry {
	return service.GetGlobalContainer().Get("orm_config").(orm.ValidatedRegistry)
}

func OrmEngineContext(ctx context.Context) *orm.Engine {
	return service.GetRequestContainer(ctx).Get("orm_context").(*orm.Engine)
}
