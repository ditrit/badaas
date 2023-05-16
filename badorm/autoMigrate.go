package badorm

import (
	"log"
	"reflect"

	"github.com/elliotchance/pie/v2"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
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

	// TODO delete
	listOrdered := db.Migrator().(postgres.Migrator).ReorderModels(listOfTables, true)
	for _, element := range listOrdered {
		log.Println(reflect.TypeOf(element).String())
	}

	err := db.AutoMigrate(listOfTables...)
	if err != nil {
		logger.Error("migration failed")
		return err
	}

	logger.Info("AutoMigration was executed successfully")
	return nil
}
