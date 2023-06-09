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

func GetCRUDServiceModule[T any]() fx.Option {
	entity := *new(T)

	moduleName := fmt.Sprintf(
		"%TCRUDServiceModule",
		entity,
	)

	kind := getBaDORMModelKind(entity)
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
			fx.Provide(NewCRUDRepository[T, uint]),
			// service
			fx.Provide(NewCRUDService[T, uint]),
		)
	default:
		log.Printf("type %T is not a BaDORM Module\n", entity)
		return fx.Invoke(failNotBadORMModule())
	}
}

func failNotBadORMModule() error {
	return fmt.Errorf("type is not a BaDORM Module")
}

func GetCRUDUnsafeServiceModule[T any]() fx.Option {
	entity := *new(T)

	moduleName := fmt.Sprintf(
		"%TCRUDUnsafeServiceModule",
		entity,
	)

	kind := getBaDORMModelKind(entity)
	switch kind {
	case KindUUIDModel:
		return fx.Module(
			moduleName,
			// models
			fx.Invoke(AddUnsafeModel[T]),
			// repository
			fx.Provide(NewCRUDUnsafeRepository[T, UUID]),
			// service
			fx.Provide(NewCRUDUnsafeService[T, UUID]),
		)
	case KindUIntModel:
		return fx.Module(
			moduleName,
			// models
			fx.Invoke(AddUnsafeModel[T]),
			// repository
			fx.Provide(NewCRUDUnsafeRepository[T, uint]),
			// service
			fx.Provide(NewCRUDUnsafeService[T, uint]),
		)
	default:
		log.Printf("type %T is not a BaDORM Module\n", entity)
		return fx.Invoke(failNotBadORMModule())
	}
}

type modelKind uint

const (
	KindUUIDModel modelKind = iota
	KindUIntModel
	KindNotBaDORMModel
)

func getBaDORMModelKind(entity any) modelKind {
	entityType := getEntityType(entity)

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

var modelsMapping = map[string]reflect.Type{}

func AddUnsafeModel[T any]() {
	entity := *new(T)
	entityType := reflect.TypeOf(entity)
	modelsMapping[entityType.Name()] = entityType
}
