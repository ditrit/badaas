package badorm

import (
	"time"

	"github.com/ditrit/badaas/persistence/gormdatabase/gormzap"
	"github.com/elliotchance/pie/v2"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetCRUD[T any, ID BadaasID](db *gorm.DB) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(db, repository), repository
}

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	allModels := pie.Flat(modelsLists)
	return db.AutoMigrate(allModels...)
}

func ConnectToDSN(
	logger *zap.Logger,
	dsn string,
	retryAmount uint,
	retryTime time.Duration,
) (*gorm.DB, error) {
	var err error
	var database *gorm.DB
	for numberRetry := uint(0); numberRetry < retryAmount; numberRetry++ {
		database, err = initializeDBFromDsn(dsn, logger)
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

// Initialize the database with the dsn string
func initializeDBFromDsn(dsn string, logger *zap.Logger) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormzap.New(logger),
	})

	if err != nil {
		return nil, err
	}

	rawDatabase, err := database.DB()
	if err != nil {
		return nil, err
	}
	// ping the underlying database
	err = rawDatabase.Ping()
	if err != nil {
		return nil, err
	}
	return database, nil
}
