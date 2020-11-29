package request

import (
	"github.com/sarulabs/di"
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/service"

	"github.com/summer-solutions/orm"
)

var OrmEngineRequestService server.InitHandler = func(s *server.Server, def *server.Def) {

	def.Name = service.OrmContextService
	def.Build = func(ctn di.Container) (interface{}, error) {
		ormConfigService := ctn.Get(service.OrmConfigService).(orm.ValidatedRegistry)
		ormEngine := ormConfigService.CreateEngine()
		ormEngine.SetLogMetaData("Source", "web-api")

		//if s.IsInLocalMode() {
		//	ormEngine.EnableQueryDebug(orm.QueryLoggerSourceDB)
		//}
		return ormEngine, nil
	}
}
