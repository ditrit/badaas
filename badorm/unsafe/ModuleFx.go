package unsafe

import (
	"fmt"
	"log"
	"reflect"

	"go.uber.org/fx"

	"github.com/ditrit/badaas/badorm"
)

func GetCRUDServiceModule[T badorm.Model]() fx.Option {
	entity := *new(T)

	moduleName := fmt.Sprintf(
		"unsafe.%TCRUDServiceModule",
		entity,
	)

	kind := badorm.GetBaDORMModelKind(entity)
	switch kind {
	case badorm.KindUUIDModel:
		return fx.Module(
			moduleName,
			// models
			fx.Invoke(addUnsafeModel[T]),
			// repository
			fx.Provide(NewCRUDRepository[T, badorm.UUID]),
			// service
			fx.Provide(NewCRUDUnsafeService[T, badorm.UUID]),
		)
	case badorm.KindUIntModel:
		return fx.Module(
			moduleName,
			// models
			fx.Invoke(addUnsafeModel[T]),
			// repository
			fx.Provide(NewCRUDRepository[T, badorm.UIntID]),
			// service
			fx.Provide(NewCRUDUnsafeService[T, badorm.UIntID]),
		)
	case badorm.KindNotBaDORMModel:
		log.Printf("type %T is not a BaDORM Module\n", entity)
		return fx.Invoke(badorm.FailNotBadORMModule())
	}

	return fx.Invoke(badorm.FailNotBadORMModule())
}

var modelsMapping = map[string]reflect.Type{}

func addUnsafeModel[T badorm.Model]() {
	entity := *new(T)
	entityType := reflect.TypeOf(entity)
	modelsMapping[entityType.Name()] = entityType
}
