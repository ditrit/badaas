package badorm

import (
	"fmt"
	"log"
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

func GetCRUDServiceModule[T Model]() fx.Option {
	entity := *new(T)

	moduleName := fmt.Sprintf(
		"%TCRUDServiceModule",
		entity,
	)

	kind := GetBaDORMModelKind(entity)
	switch kind {
	case KindUUIDModel:
		return fx.Module(
			moduleName,
			// repository
			fx.Provide(NewCRUDRepository[T, UUID]),
			// service
			fx.Provide(NewCRUDService[T, UUID]),
		)
	case KindUIntModel:
		return fx.Module(
			moduleName,
			// repository
			fx.Provide(NewCRUDRepository[T, UIntID]),
			// service
			fx.Provide(NewCRUDService[T, UIntID]),
		)
	case KindNotBaDORMModel:
		log.Printf("type %T is not a BaDORM Module\n", entity)
		return fx.Invoke(FailNotBadORMModule())
	}

	return fx.Invoke(FailNotBadORMModule())
}

func FailNotBadORMModule() error {
	return fmt.Errorf("type is not a BaDORM Module")
}

type ModelKind uint

const (
	KindUUIDModel ModelKind = iota
	KindUIntModel
	KindNotBaDORMModel
)

func GetBaDORMModelKind(entity Model) ModelKind {
	entityType := GetEntityType(entity)

	_, isUUIDModel := entityType.FieldByName("UUIDModel")
	if isUUIDModel {
		return KindUUIDModel
	}

	_, isUIntModel := entityType.FieldByName("UIntModel")
	if isUIntModel {
		return KindUIntModel
	}

	return KindNotBaDORMModel
}

// Get the reflect.Type of any entity or pointer to entity
func GetEntityType(entity any) reflect.Type {
	entityType := reflect.TypeOf(entity)

	// entityType will be a pointer if the relation can be nullable
	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
	}

	return entityType
}
