package gormdatabase

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/configuration"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Create the dsn string from the configuration
func createDialectorFromConf(databaseConfiguration configuration.DatabaseConfiguration) gorm.Dialector {
	switch databaseConfiguration.GetDialector() {
	case configuration.PostgreSQL:
		return badorm.CreatePostgreSQLDialector(
			databaseConfiguration.GetHost(),
			databaseConfiguration.GetUsername(),
			databaseConfiguration.GetPassword(),
			databaseConfiguration.GetSSLMode(),
			databaseConfiguration.GetDBName(),
			databaseConfiguration.GetPort(),
		)
	case configuration.MySQL:
		return badorm.CreateMySQLDialector(
			databaseConfiguration.GetHost(),
			databaseConfiguration.GetUsername(),
			databaseConfiguration.GetPassword(),
			databaseConfiguration.GetSSLMode(),
			databaseConfiguration.GetDBName(),
			databaseConfiguration.GetPort(),
		)
	}

	return nil
}

// Creates the database object with using the database configuration and exec the setup
func SetupDatabaseConnection(
	logger *zap.Logger,
	databaseConfiguration configuration.DatabaseConfiguration,
) (*gorm.DB, error) {
	return badorm.ConnectToDialector(
		logger,
		createDialectorFromConf(databaseConfiguration),
		databaseConfiguration.GetRetry(),
		databaseConfiguration.GetRetryTime(),
	)
}
