package logger

import (
	"go.uber.org/zap"
)

type LoggerType int

const (
	_ = iota

	//
	ProductionLogger
	DevelopmentLogger
)

// Initialize zap global logger instance
func InitLogger(loggerType LoggerType) error {
	var logger *zap.Logger
	var err error
	switch loggerType {
	case ProductionLogger:
		logger, err = zap.NewProduction()
		if err != nil {
			return err
		}
	case DevelopmentLogger:
		logger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	}
	zap.ReplaceGlobals(logger)
	return nil
}
