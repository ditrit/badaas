package badorm

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/persistence/gormdatabase/gormzap"
)

func CreatePostgreSQLDialector(host, username, password, sslmode, dbname string, port int) gorm.Dialector {
	return postgres.Open(CreatePostgreSQLDSN(
		host, username, password, sslmode, dbname, port,
	))
}

func CreatePostgreSQLDSN(host, username, password, sslmode, dbname string, port int) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
		username, password, host, port, sslmode, dbname,
	)
}

func CreateMySQLDialector(host, username, password, dbname string, port int) gorm.Dialector {
	return mysql.Open(CreateMySQLDSN(
		host, username, password, dbname, port,
	))
}

func CreateMySQLDSN(host, username, password, dbname string, port int) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname,
	)
}

func CreateSQLiteDialector(path string) gorm.Dialector {
	return sqlite.Open(CreateSQLiteDSN(path))
}

func CreateSQLiteDSN(path string) string {
	return fmt.Sprintf("sqlite:%s", path)
}

func ConnectToDialector(
	logger *zap.Logger,
	dialector gorm.Dialector,
	retryAmount uint,
	retryTime time.Duration,
) (*gorm.DB, error) {
	var err error
	var database *gorm.DB
	for numberRetry := uint(0); numberRetry < retryAmount; numberRetry++ {
		database, err = gorm.Open(dialector, &gorm.Config{
			Logger: gormzap.New(logger),
		})

		if err == nil {
			logger.Sugar().Debugf("Database connection is active")
			return database, nil
		}

		logger.Sugar().Debugf("Database connection failed with error %q", err.Error())
		logger.Sugar().Debugf(
			"Retrying database connection %d/%d in %s",
			numberRetry+1, retryAmount, retryTime.String(),
		)
		time.Sleep(retryTime)
	}

	return nil, err
}
