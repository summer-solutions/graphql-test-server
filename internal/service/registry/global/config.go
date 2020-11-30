package global

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/internal/service/config"

	"github.com/sarulabs/di"
)

var ConfigGlobalService server.InitHandler = func(s *server.Server, def *server.Def) {
	def.Name = "config"
	def.Build = func(ctn di.Container) (interface{}, error) {
		return config.NewViperConfig("../../config/web-api/config.local.yaml")
	}
}
