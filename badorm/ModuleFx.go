package badorm

import (
	"fmt"
	"reflect"

	"go.uber.org/fx"
)

type GetModelsResult struct {
	fx.Out

	Models []any `group:"modelsTables"`
}

var BaDORMModule = fx.Module(
	"BaDORM",
	fx.Invoke(
		fx.Annotate(
			autoMigrate,
			fx.ParamTags(`group:"modelsTables"`),
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
		fx.Invoke(AddModel[T]),
		// repository
		fx.Provide(NewCRUDRepository[T, ID]),
		// service
		fx.Provide(NewCRUDService[T, ID]),
	)
}

var modelsMapping = map[string]reflect.Type{}

func AddModel[T any]() {
	entity := *new(T)
	entityType := reflect.TypeOf(entity)
	modelsMapping[entityType.Name()] = entityType
}