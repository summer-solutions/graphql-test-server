package global

import (
	"summer-solutions/graphql-test-server/internal/server"

	"github.com/apex/log"
	"github.com/sarulabs/di"
)

var LogGlobalService server.InitHandler = func(s *server.Server, def *server.Def) {
	def.Name = "log"
	def.Build = func(ctn di.Container) (interface{}, error) {
		return log.WithFields(&log.Fields{}), nil
	}
}
