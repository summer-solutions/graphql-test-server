package service

import (
	"context"
	log2 "github.com/apex/log"
	"summer-solutions/graphql-test-server/internal/log"
)

//context services
const OrmContextService string = "orm_context"

//global services
const ConfigService string = "config"
const OrmConfigService string = "orm_config"

func Log(context context.Context) log2.Interface {
	return log.FromContext(context)
}
