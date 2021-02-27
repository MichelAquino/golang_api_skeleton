package context

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is an interface to log object.
type Logger interface {
	// Debug writes a debug log level.
	Debug(class, method, requestID, ip, action, result, message string)

	// Info writes an info log level.
	Info(class, method, requestID, ip, action, result, message string)

	// Warn writes a warning log level.
	Warn(class, method, requestID, ip, action, result, message string)

	// Error writes an error log level.
	Error(class, method, requestID, ip, action, result, message string)
}

// APILog is the API logger.
type APILog struct {
	logger *zap.Logger
}

// NewAPILog returns a pointer of the APILog.
func NewAPILog() *APILog {
	loggerInstance := getNewLogInstance()
	return &APILog{
		logger: loggerInstance,
	}
}

// Debug writes a debug log level.
func (l APILog) Debug(class, method, requestID, ip, action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(class, method, requestID, ip, action, result)
	l.logger.Debug(message, fields...)
}

// Info writes an info log level.
func (l APILog) Info(class, method, requestID, ip, action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(class, method, requestID, ip, action, result)
	l.logger.Info(message, fields...)
}

// Warn writes a warning log level.
func (l APILog) Warn(class, method, requestID, ip, action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(class, method, requestID, ip, action, result)
	l.logger.Warn(message, fields...)
}

// Error writes an error log level.
func (l APILog) Error(class, method, requestID, ip, action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(class, method, requestID, ip, action, result)
	l.logger.Error(message, fields...)
}

// getLogFields return a logrus fields instance.
func getLogFields(class, method, requestID, ip, action, result string) []zapcore.Field {
	return []zapcore.Field{
		zap.String("request_id", requestID),
		zap.String("struct", class),
		zap.String("method", method),
		zap.String("ip", ip),
		zap.String("action", action),
		zap.String("result", result),
	}
}

var apiLogger *APILog
var onceLog sync.Once

// GetLogger returns a new instance of the log.
func GetLogger() Logger {
	onceLog.Do(func() {
		apiLogger = NewAPILog()
	})

	return apiLogger
}

// getNewLogInstance returns a new instance of Logrus log.
func getNewLogInstance() *zap.Logger {
	logLevel := getLogLevel()

	logOutputPaths := getLogOutputPaths()

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		// EncoderConfig: zap.NewProductionEncoderConfig(),
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      logOutputPaths,
		ErrorOutputPaths: logOutputPaths,
	}

	logger, _ := config.Build()
	return logger
}

func getLogLevel() zapcore.Level {
	logLevelConfig := strings.ToLower(viper.GetString("log.level"))
	if logLevelConfig == "debug" {
		return zap.DebugLevel
	}

	if logLevelConfig == "info" {
		return zap.InfoLevel
	}

	if logLevelConfig == "warn" {
		return zap.WarnLevel
	}

	return zap.ErrorLevel
}

func getLogOutputPaths() []string {
	paths := []string{"stderr"}

	sendLogToFile := viper.GetBool("log.to.file")
	if sendLogToFile {
		logFileName := viper.GetString("log.file.name")
		paths = append(paths, logFileName)
	}

	return paths
}
