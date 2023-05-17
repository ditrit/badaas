package badorm

import (
	"reflect"

	"github.com/elliotchance/pie/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func autoMigrate(modelsLists [][]any, db *gorm.DB, logger *zap.Logger) error {
	allModels := pie.Flat(modelsLists)
	registeredModels := pie.Map(allModels, func(model any) string {
		return reflect.TypeOf(model).String()
	})
	logger.Sugar().Debug(
		"Registered models: ",
		registeredModels,
	)

	err := db.AutoMigrate(allModels...)
	if err != nil {
		logger.Error("migration failed")
		return err
	}

	logger.Info("AutoMigration was executed successfully")
	return nil
}
