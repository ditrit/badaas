package logger_test

import (
	"testing"

	"github.com/ditrit/badaas/logger"
)

func TestInitializeDevelopmentLogger(t *testing.T) {
	err := logger.InitLogger(logger.DevelopmentLogger)
	if err != nil {
		t.Errorf("InitLogger should return a null value")
	}
}

func TestInitializeProductionLogger(t *testing.T) {
	err := logger.InitLogger(logger.ProductionLogger)
	if err != nil {
		t.Errorf("InitLogger should return a null value")
	}
}
