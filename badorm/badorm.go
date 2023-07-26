package badorm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm/logger"
)

func GetCRUD[T Model, ID ModelID](
	logger logger.Interface,
	db *gorm.DB,
) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(logger, db, repository), repository
}

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	allModels := pie.Flat(modelsLists)
	return db.AutoMigrate(allModels...)
}
