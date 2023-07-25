package configuration

import (
	"github.com/ditrit/badaas/configuration/defaults"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func NewLoggerCommandsInitializer() CommandsInitializer {
	modeValidator := validators.AuthorizedValues(ProductionLogger, DevelopmentLogger)

	return CommandsKeyInitializer{
		KeySetter: NewKeySetter(),
		Keys: []Key{
			{
				Name:     LoggerRequestTemplateKey,
				ValType:  verdeter.IsStr,
				Usage:    "Template message for all request logs",
				DefaultV: "Receive {{method}} request on {{url}}",
			},
			{
				Name:      LoggerModeKey,
				ValType:   verdeter.IsStr,
				Usage:     "The logger mode (default to \"prod\")",
				DefaultV:  ProductionLogger,
				Validator: &modeValidator,
			},
			{
				Name:     LoggerDisableStacktraceKey,
				ValType:  verdeter.IsBool,
				Usage:    "Disable error stacktrace from logs (default to true)",
				DefaultV: true,
			},
			{
				Name:     LoggerSlowThresholdKey,
				ValType:  verdeter.IsInt,
				Usage:    "Threshold for the slow query warning in milliseconds (default to 200)",
				DefaultV: defaults.LoggerSlowThreshold,
			},
			{
				Name:     LoggerIgnoreRecordNotFoundErrorKey,
				ValType:  verdeter.IsBool,
				Usage:    "If true, ignore gorm.ErrRecordNotFound error for logger (default to false)",
				DefaultV: false,
			},
			{
				Name:     LoggerParameterizedQueriesKey,
				ValType:  verdeter.IsBool,
				Usage:    "If true, don't include params in the query execution logs (default to false)",
				DefaultV: false,
			},
		},
	}
}
