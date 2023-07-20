package gormdatabase

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/gormdatabase/gormzap"
)

func createDialectorFromConf(databaseConfiguration configuration.DatabaseConfiguration) (gorm.Dialector, error) {
	switch databaseConfiguration.GetDialector() {
	case configuration.PostgreSQL:
		return badorm.CreatePostgreSQLDialector(
			databaseConfiguration.GetHost(),
			databaseConfiguration.GetUsername(),
			databaseConfiguration.GetPassword(),
			databaseConfiguration.GetSSLMode(),
			databaseConfiguration.GetDBName(),
			databaseConfiguration.GetPort(),
		), nil
	case configuration.MySQL:
		return badorm.CreateMySQLDialector(
			databaseConfiguration.GetHost(),
			databaseConfiguration.GetUsername(),
			databaseConfiguration.GetPassword(),
			databaseConfiguration.GetDBName(),
			databaseConfiguration.GetPort(),
		), nil
	case configuration.SQLite:
		return badorm.CreateSQLiteDialector(
			databaseConfiguration.GetHost(),
		), nil
	case configuration.SQLServer:
		return badorm.CreateSQLServerDialector(
			databaseConfiguration.GetHost(),
			databaseConfiguration.GetUsername(),
			databaseConfiguration.GetPassword(),
			databaseConfiguration.GetDBName(),
			databaseConfiguration.GetPort(),
		), nil
	default:
		return nil, fmt.Errorf("unknown dialector: %s", databaseConfiguration.GetDialector())
	}
}

// TODO cambiar el nombre de este paquete y actualizar el file struct

// Creates the database object with using the database configuration and exec the setup
func SetupDatabaseConnection(
	logger *zap.Logger,
	databaseConfiguration configuration.DatabaseConfiguration,
) (*gorm.DB, error) {
	dialector, err := createDialectorFromConf(databaseConfiguration)
	if err != nil {
		return nil, err
	}

	return badorm.ConnectToDialector(
		// TODO aca me gustaria que se pueda configurar todo el gormzap tambien
		gormzap.NewDefault(logger),
		dialector,
		databaseConfiguration.GetRetry(),
		databaseConfiguration.GetRetryTime(),
	)
}
