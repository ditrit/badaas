package badorm

import (
	"reflect"

	"gorm.io/gorm"
)

func GetCRUD[T any, ID BadaasID](db *gorm.DB) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	AddModel[T]()
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(db, repository), repository
}

// TODO verificar si esto va a seguir siendo util o no
var modelsMapping = map[string]reflect.Type{}

// TODO no deberia ser exportado
func AddModel[T any]() {
	entity := *new(T)
	entityType := reflect.TypeOf(entity)
	modelsMapping[entityType.Name()] = entityType
}
