package badorm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

func GetCRUD[T Model, ID ModelID](db *gorm.DB) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(db, repository), repository
}

// TODO auto migracion no obligatoria
func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	allModels := pie.Flat(modelsLists)
	return db.AutoMigrate(allModels...)
}
