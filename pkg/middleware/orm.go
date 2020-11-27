package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/summer-solutions/orm"
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/service"
)

const (
	SourceWebAPI = "web-api"
)

func Orm(source string) func(*gin.Engine) error {
	return func(engine *gin.Engine) error {
		engine.Use(func(c *gin.Context) {
			ormConfigService := server.IoC.Get(service.OrmConfigService).(orm.ValidatedRegistry)
			ormEngine := ormConfigService.CreateEngine()
			ormEngine.SetLogMetaData("IP", c.ClientIP())

			c.Set(service.OrmContextService, ormEngine)
			ormEngine.SetLogMetaData("Source", source)

			//if s.IsInLocalMode() {
			//	ormEngine.EnableQueryDebug(orm.QueryLoggerSourceDB)
			//}
		})

		return nil
	}
}
