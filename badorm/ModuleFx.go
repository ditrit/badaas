package badorm

import (
	"fmt"
	"reflect"

	"go.uber.org/fx"
)

var BaDORMModule = fx.Module(
	"BaDORM",
	fx.Invoke(
		fx.Annotate(
			autoMigrate,
			fx.ParamTags(`group:"modelTables"`),
		),
	),
)

func GetCRUDServiceModule[T any, ID BadaasID]() fx.Option {
	return fx.Module(
		fmt.Sprintf(
			"%TCRUDServiceModule",
			*new(T),
		),
		// models
		fx.Provide(AddModel[T]),
		// repository
		fx.Provide(NewCRUDRepository[T, ID]),
		// service
		fx.Provide(NewCRUDService[T, ID]),
	)
}

var modelsMapping = map[string]reflect.Type{}

type AddModelResult struct {
	fx.Out

	Model any `group:"modelTables"`
}

func AddModel[T any]() AddModelResult {
	entity := *new(T)
	entityType := reflect.TypeOf(entity)
	modelsMapping[entityType.Name()] = entityType

	return AddModelResult{
		Model: entity,
	}
}
