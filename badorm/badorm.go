package badorm

import (
	"github.com/elliotchance/pie/v2"
)

func GetCRUD[T Model, ID ModelID](
	db *DB,
) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(db, repository), repository
}

func autoMigrate(modelsLists [][]any, db *DB) error {
	allModels := pie.Flat(modelsLists)
	return db.GormDB.AutoMigrate(allModels...)
}
