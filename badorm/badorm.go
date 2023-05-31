package badorm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

func GetCRUD[T any, ID BadaasID](db *gorm.DB) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(db, repository), repository
}

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	allModels := pie.Flat(modelsLists)
	return db.AutoMigrate(allModels...)
}
