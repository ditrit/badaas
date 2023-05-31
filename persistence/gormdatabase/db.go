package gormdatabase

import (
	"fmt"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/configuration"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Create the dsn string from the configuration
func createDsnFromConf(databaseConfiguration configuration.DatabaseConfiguration) string {
	dsn := createDsn(
		databaseConfiguration.GetHost(),
		databaseConfiguration.GetUsername(),
		databaseConfiguration.GetPassword(),
		databaseConfiguration.GetSSLMode(),
		databaseConfiguration.GetDBName(),
		databaseConfiguration.GetPort(),
	)
	return dsn
}

// Create the dsn strings with the provided args
func createDsn(host, username, password, sslmode, dbname string, port int) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
		username, password, host, port, sslmode, dbname,
	)
}

// Creates the database object with using the database configuration and exec the setup
func SetupDatabaseConnection(
	logger *zap.Logger,
	databaseConfiguration configuration.DatabaseConfiguration,
) (*gorm.DB, error) {
	dsn := createDsnFromConf(databaseConfiguration)
	return badorm.ConnectToDSN(
		logger,
		dsn,
		databaseConfiguration.GetRetry(),
		databaseConfiguration.GetRetryTime(),
	)
}
