package request

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/service"

	"github.com/sarulabs/di"
)

var OrmEngineRequestService server.InitHandler = func(s *server.Server, def *server.Def) {

	def.Name = "orm_engine"
	def.Build = func(ctn di.Container) (interface{}, error) {
		ormConfigService := service.OrmConfig()
		ormEngine := ormConfigService.CreateEngine()
		ormEngine.SetLogMetaData("Source", "web-api")

		//if s.IsInLocalMode() {
		//	ormEngine.EnableQueryDebug(orm.QueryLoggerSourceDB)
		//}
		return ormEngine, nil
	}
}
