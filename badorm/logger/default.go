package logger

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

var (
	Default       = New(DefaultConfig)
	DefaultWriter = WriterWrapper{Writer: log.New(os.Stdout, "\r\n", log.LstdFlags)}
)

type defaultLogger struct {
	gormLogger.Interface
	Config
}

func New(config Config) Interface {
	return &defaultLogger{
		Config:    config,
		Interface: gormLogger.New(DefaultWriter, config.toGormConfig()),
	}
}

func (l *defaultLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l.ToLogMode(level)
}

func (l *defaultLogger) ToLogMode(level gormLogger.LogLevel) Interface {
	newLogger := *l
	newLogger.LogLevel = level
	newLogger.Interface = newLogger.Interface.LogMode(level)

	return &newLogger
}

const nanoToMicro = 1e6

func (l defaultLogger) TraceTransaction(ctx context.Context, begin time.Time) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case l.SlowTransactionThreshold != DisableThreshold && elapsed > l.SlowTransactionThreshold && l.LogLevel >= gormLogger.Warn:
		l.Interface.Warn(ctx, "transaction_slow (>= %vms) [%.3fms]", l.SlowTransactionThreshold, float64(elapsed.Nanoseconds())/nanoToMicro)
	case l.LogLevel >= gormLogger.Info:
		l.Interface.Info(ctx, "transaction_exec [%.3fms]", float64(elapsed.Nanoseconds())/nanoToMicro)
	}
}

type WriterWrapper struct {
	Writer gormLogger.Writer
}

// Info, Warn, Error or Trace + Printf
const defaultStacktraceLen = 2

func (w WriterWrapper) Printf(msg string, args ...interface{}) {
	if len(args) > 0 {
		// change the file path to avoid showing badorm internal files
		firstArg := args[0]

		_, isString := firstArg.(string)
		if isString {
			file, line, caller := FindLastCaller(defaultStacktraceLen)
			if caller != 0 {
				w.Writer.Printf(
					msg,
					append(
						[]any{file + ":" + strconv.FormatInt(int64(line), 10)},
						args[1:]...,
					)...,
				)

				return
			}
		}
	}

	w.Writer.Printf(msg, args...)
}
