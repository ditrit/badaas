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
	// TODO verificar que sea un badorm module?
	// TODO sacar solo cual es el id?
	return fx.Module(
		fmt.Sprintf(
			"%TCRUDServiceModule",
			*new(T),
		),
		// repository
		fx.Provide(NewCRUDRepository[T, ID]),
		// service
		fx.Provide(NewCRUDService[T, ID]),
	)
}

func GetCRUDUnsafeServiceModule[T any, ID BadaasID]() fx.Option {
	// TODO verificar que sea un badorm module?
	// TODO sacar solo cual es el id?
	return fx.Module(
		fmt.Sprintf(
			"%TCRUDUnsafeServiceModule",
			*new(T),
		),
		// models
		fx.Invoke(AddUnsafeModel[T]),
		// repository
		fx.Provide(NewCRUDUnsafeRepository[T, ID]),
		// service
		fx.Provide(NewCRUDUnsafeService[T, ID]),
	)
}

var modelsMapping = map[string]reflect.Type{}

func AddUnsafeModel[T any]() {
	entity := *new(T)
	entityType := reflect.TypeOf(entity)
	modelsMapping[entityType.Name()] = entityType
}
