package logger

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

const (
	// Silent silent log level
	Silent gormLogger.LogLevel = gormLogger.Silent
	// Error error log level
	Error gormLogger.LogLevel = gormLogger.Error
	// Warn warn log level
	Warn gormLogger.LogLevel = gormLogger.Warn
	// Info info log level
	Info gormLogger.LogLevel = gormLogger.Info
)

var (
	Default = New(gormLogger.Config{
		SlowThreshold:             200 * time.Millisecond, //nolint:gomnd // default definition
		LogLevel:                  gormLogger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	DefaultWriter = WriterWrapper{Writer: log.New(os.Stdout, "\r\n", log.LstdFlags)}
)

func New(config gormLogger.Config) gormLogger.Interface {
	return gormLogger.New(DefaultWriter, config)
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

var (
	gormPackage   = filepath.Join("gorm.io", "gorm")
	badormPackage = filepath.Join("badaas", "badorm")
)

// search in the stacktrace the last file outside gormzap, badorm and gorm
func FindLastCaller(skip int) (string, int, int) {
	// +1 because at least one will be inside gorm
	// +1 because of this function
	for i := skip + 1 + 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)

		if !ok {
			// we checked in all the stacktrace and none meet the conditions,
			return "", 0, 0
		} else if (!strings.Contains(file, gormPackage) && !strings.Contains(file, badormPackage)) || strings.HasSuffix(file, "_test.go") {
			// file outside badorm and gorm or a test file (util for badorm developers)
			return file, line, i
		}
	}
}
