package badorm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

func AutoMigrate(models []any, db *gorm.DB) error {
	return autoMigrate([][]any{models}, db)
}

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	allModels := pie.Flat(modelsLists)

	err := db.AutoMigrate(allModels...)
	if err != nil {
		return err
	}

	return nil
}
