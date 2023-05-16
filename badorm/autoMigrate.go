package badorm

import (
	"reflect"

	"github.com/elliotchance/pie/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func autoMigrate(listOfTables []any, db *gorm.DB, logger *zap.Logger) error {
	registeredModels := pie.Map(listOfTables, func(model any) string {
		return reflect.TypeOf(model).String()
	})
	logger.Sugar().Debug(
		"Registered models: ",
		registeredModels,
	)

	err := db.AutoMigrate(listOfTables...)
	if err != nil {
		logger.Error("migration failed")
		return err
	}

	logger.Info("AutoMigration was executed successfully")
	return nil
}
