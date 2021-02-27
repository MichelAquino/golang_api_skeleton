package log

import (
	"context"
	"sync"

	"github.com/spf13/viper"
)

// Logger is a interface to log object
type Logger interface {
	// Debug write a debug log level
	Debug(context context.Context, action string, message string, extraFields map[string]string)

	// Info write a info log level
	Info(context context.Context, action string, message string, extraFields map[string]string)

	// Warn write a warning log level
	Warn(context context.Context, action string, message string, extraFields map[string]string)

	// Error write a error log level
	Error(context context.Context, action string, message string, extraFields map[string]string)

	// Fatal write a fatal log level. The logger then calls os.Exit(1), even if logging at Fatal level is disabled.
	Fatal(context context.Context, action string, message string, extraFields map[string]string)

	// Panic write a Panic log level. Then panics then panic, even if logging at Panic level is disabled.
	Panic(context context.Context, action string, message string, extraFields map[string]string)
}

var appLog *zapLog
var onceLog sync.Once

// GetLogger return a new instance of the log for the application
func GetLogger() Logger {
	onceLog.Do(func() {
		encoding := viper.GetString("log.encoding")
		if encoding == "" {
			encoding = jsonEncoding
		}

		appLog = newZapLog(encoding)
	})

	return appLog
}
