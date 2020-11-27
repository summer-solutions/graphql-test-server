package initializer

import (
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/service"
	"summer-solutions/graphql-test-server/pkg/service/config"

	"github.com/sarulabs/di"
)

var ConfigHandler server.InitHandler = func(s *server.Spring) error {

	return server.IoCBuilder.Add(
		di.Def{
			Name: service.ConfigService,
			Build: func(ctn di.Container) (interface{}, error) {
				return config.NewViperConfig("../../config/web-api/config.local.yaml")
			},
		},
	)
}
