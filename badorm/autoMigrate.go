package badorm

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func autoMigrate(listOfTables []any, db *gorm.DB, logger *zap.Logger) error {
	err := db.AutoMigrate(listOfTables...)
	if err != nil {
		logger.Error("migration failed")
		return err
	}

	logger.Info("AutoMigration was executed successfully")
	return nil
}
