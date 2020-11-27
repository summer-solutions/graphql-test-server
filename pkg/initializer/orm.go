package initializer

import (
	"io/ioutil"
	"os"
	"summer-solutions/graphql-test-server/internal/server"
	"summer-solutions/graphql-test-server/pkg/entity"
	"summer-solutions/graphql-test-server/pkg/service"

	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/sarulabs/di"
	"github.com/summer-solutions/orm"
)

var ormConfig orm.ValidatedRegistry

var OrmHandler server.InitHandler = func(s *server.Spring) error {
	registry, err := initOrmRegistry(s)
	if err != nil {
		panic(err)
	}

	registerEntities(registry)

	err = initOrmConfig(registry)
	if err != nil {
		panic(err)
	}

	if s.IsInLocalMode() {
		err = checkForAlters(ormConfig)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func registerEntities(registry *orm.Registry) {
	var UserEntity *entity.UserEntity

	//entities
	registry.RegisterEntity(UserEntity)

	//enums
	//registry.RegisterEnumStruct("entity.DurationRangeAll", entity.DurationRangeAll)
}

func initOrmRegistry(_ *server.Spring) (*orm.Registry, error) {
	var yamlFileData []byte
	var err error

	if os.Getenv("ORM_CONFIG_FILE") != "" {
		yamlFileData, err = ioutil.ReadFile(os.Getenv("ORM_CONFIG_FILE"))
	} else {
		yamlFileData, err = ioutil.ReadFile("../../config/orm/config.yaml")
	}

	if err != nil {
		return nil, err
	}

	var parsedYaml map[string]interface{}

	err = yaml.Unmarshal(yamlFileData, &parsedYaml)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	for k, v := range parsedYaml["orm"].(map[interface{}]interface{}) {
		data[k.(string)] = v
	}

	return orm.InitByYaml(data), nil
}

func initOrmConfig(registry *orm.Registry) error {
	var err error
	ormConfig, err = registry.Validate()
	if err != nil {
		return err
	}
	return server.IoCBuilder.Add(
		di.Def{
			Name: service.OrmConfigService,
			Build: func(ctn di.Container) (interface{}, error) {
				return ormConfig, nil
			},
		},
	)
}

func checkForAlters(ormConfig orm.ValidatedRegistry) error {
	engine := ormConfig.CreateEngine()

	alters := engine.GetAlters()
	for _, alter := range alters {
		if alter.Safe {
			color.Green("%s\n\n", alter.SQL)
		} else {
			color.Red("%s\n\n", alter.SQL)
		}
	}

	return nil
}
