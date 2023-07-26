package badorm

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm/logger"
)

type Config struct {
	gorm.Config
	// Logger for badorm and gorm
	Logger logger.Interface
	// amount of times to retry db connection
	RetryAmount uint
	// time between retries db connection
	RetryTime time.Duration
}

const defaultRetryTime = time.Duration(5 * time.Second)

func (c Config) toGormConfig() *gorm.Config {
	gormConfig := c.Config
	gormConfig.Logger = c.Logger

	return &gormConfig
}

func (c Config) getRetry() (uint, time.Duration) {
	retryAmount := c.RetryAmount
	if retryAmount == 0 {
		// try at least once
		retryAmount = 1
	}

	retryTime := c.RetryTime
	if retryTime == 0 {
		retryTime = defaultRetryTime
	}

	return retryAmount, retryTime
}

type DB struct {
	GormDB *gorm.DB
	Logger logger.Interface
}

// Open a new db connection using the "dialector" and,
// if connection is established, config the DB object with "config"
// Check Config.RetryAmount and Config.RetryTime for retrying
func Open(dialector gorm.Dialector, config Config) (*DB, error) {
	var err error

	retryAmount, retryTime := config.getRetry()
	logger := config.Logger

	for retryNumber := uint(0); retryNumber < retryAmount; retryNumber++ {
		gormDB, err := gorm.Open(dialector, config.toGormConfig())

		if err == nil {
			logger.Info(context.Background(), "Database connection is active")

			return &DB{
				GormDB: gormDB,
				Logger: logger,
			}, nil
		}

		// there are more retries
		if retryNumber < retryAmount-1 {
			logger.Info(
				context.Background(),
				"Database connection failed with error %q, retrying %d/%d in %s",
				err.Error(),
				retryNumber+1+1, // +1 for counting from 1 and +1 for next iteration
				retryAmount,
				retryTime,
			)
			time.Sleep(retryTime)
		}
	}

	logger.Error(
		context.Background(),
		"Database connection failed with error %q",
		err.Error(),
	)

	return nil, err
}

func New(gormDB *gorm.DB, config Config) *DB {
	return &DB{
		GormDB: gormDB,
		Logger: config.Logger,
	}
}
