package gormzap

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// TODO algo para saber cuantas queries hiciste por transaction y cuanto tiempo tardó
// en prod cuantas queries no me interesa porque es la misma que en dev
// lo que si me va a interesar son warnings de transacciones lentas,
// que no es lo mismo que solo una query lenta (aunque posiblemente tengan relacion)

type Config struct {
	SlowThreshold time.Duration // Slow SQL threshold
	// TODO ver si el zaplogger no tiene ya un logLevel tambien
	LogLevel                  gormLogger.LogLevel // Log level
	IgnoreRecordNotFoundError bool                // if true, ignore gorm.ErrRecordNotFound error for logger
	ParameterizedQueries      bool                // if true, don't include params in the SQL log
}

const IgnoreSlowQueries = 0

var DefaultConfig = Config{
	LogLevel:                  gormLogger.Warn,
	SlowThreshold:             200 * time.Millisecond,
	IgnoreRecordNotFoundError: false,
	ParameterizedQueries:      false, // TODO usar para algo
}

// This type implement the [gorm.io/gorm/logger.Interface] interface.
// It is to be used as a replacement for the original logger
type Logger struct {
	Config
	ZapLogger *zap.Logger
}

// The constructor of the gormzap logger with default config
func NewDefault(zapLogger *zap.Logger) gormLogger.Interface {
	return &Logger{
		ZapLogger: zapLogger,
		Config:    DefaultConfig,
	}
}

// The constructor of the gormzap logger
func New(zapLogger *zap.Logger, config Config) gormLogger.Interface {
	return &Logger{
		ZapLogger: zapLogger,
		Config:    config,
	}
}

// Set the log mode to the value passed as argument
func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

// log info
func (l Logger) Info(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		l.logger().Sugar().Debugf(str, args...)
	}
}

// log warning
func (l Logger) Warn(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.logger().Sugar().Warnf(str, args...)
	}
}

// log an error
func (l Logger) Error(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.logger().Sugar().Errorf(str, args...)
	}
}

// log a trace
func (l Logger) Trace(
	_ context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsedTime := time.Since(begin)

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rowsAffected := fc()
		l.logger().Error(
			"query_error",
			append(getZapFields(elapsedTime, rowsAffected, sql), zap.Error(err))...,
		)
	case l.SlowThreshold != IgnoreSlowQueries && elapsedTime > l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		sql, rowsAffected := fc()
		l.logger().Warn(
			"query_slow",
			getZapFields(elapsedTime, rowsAffected, sql)...,
		)
	case l.LogLevel >= gormLogger.Info:
		sql, rowsAffected := fc()
		l.logger().Debug(
			"query_exec",
			getZapFields(elapsedTime, rowsAffected, sql)...,
		)
	}
}

func getZapFields(elapsedTime time.Duration, rowsAffected int64, sql string) []zapcore.Field {
	rowsAffectedString := strconv.FormatInt(rowsAffected, 10)
	if rowsAffected == -1 {
		rowsAffectedString = "-"
	}

	return []zapcore.Field{
		zap.Duration("elapsed_time", elapsedTime),
		zap.String("rows_affected", rowsAffectedString),
		zap.String("sql", sql),
	}
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
	// TODO ojo con eso si lo muevo a badorm
	zapgormPackage = filepath.Join("github.com", "ditrit", "badaas", "persistence", "gormdatabase", "gormzap")
)

// TODO ver esto que mierda es
// return a logger that log the right caller
func (l Logger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)

		if ok && !strings.HasSuffix(file, "_test.go") && !strings.Contains(file, gormPackage) && !strings.Contains(file, zapgormPackage) {
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(i))
		}
	}

	return l.ZapLogger
}
