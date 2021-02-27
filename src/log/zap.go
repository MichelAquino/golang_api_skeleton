package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	appContext "github.com/michelaquino/golang_api_skeleton/src/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	keyValueEncoding string = "key-value"
	jsonEncoding     string = "json"
)

// zapLog is the API logger
type zapLog struct {
	logger *zap.Logger
}

func newZapLog(encoding string) *zapLog {
	err := zap.RegisterEncoder(keyValueEncoding, func(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewKeyValueEncoder(config), nil
	})

	if err != nil && encoding == keyValueEncoding {
		encoding = jsonEncoding
	}

	config := getConfig(encoding)
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Sprintf("can't initialize zap logger: %s", err.Error()))
	}

	return &zapLog{
		logger: logger,
	}
}

// Debug write a debug log level
func (z zapLog) Debug(ctx context.Context, action string, message string, extraFields map[string]string) {
	defer z.logger.Sync()
	fields := z.getExtraFields(ctx, action, extraFields)
	z.logger.Debug(message, fields...)
}

// Info write a info log level
func (z zapLog) Info(ctx context.Context, action string, message string, extraFields map[string]string) {
	defer z.logger.Sync()
	fields := z.getExtraFields(ctx, action, extraFields)
	z.logger.Info(message, fields...)
}

// Warn write a warning log level
func (z zapLog) Warn(ctx context.Context, action string, message string, extraFields map[string]string) {
	defer z.logger.Sync()
	fields := z.getExtraFields(ctx, action, extraFields)
	z.logger.Warn(message, fields...)
}

// Error write a error log level
func (z zapLog) Error(ctx context.Context, action string, message string, extraFields map[string]string) {
	defer z.logger.Sync()
	fields := z.getExtraFields(ctx, action, extraFields)
	z.logger.Error(message, fields...)
}

// Fatal write a fatal log level. The logger then calls os.Exit(1), even if logging at Fatal level is disabled.
func (z zapLog) Fatal(ctx context.Context, action string, message string, extraFields map[string]string) {
	defer z.logger.Sync()
	fields := z.getExtraFields(ctx, action, extraFields)
	z.logger.Fatal(message, fields...)
}

// Panic write a Panic log level. Then panics then panic, even if logging at Panic level is disabled.
func (z zapLog) Panic(ctx context.Context, action string, message string, extraFields map[string]string) {
	defer z.logger.Sync()
	fields := z.getExtraFields(ctx, action, extraFields)
	z.logger.Panic(message, fields...)
}

// getExtraFields return a zap fields instance
func (z zapLog) getExtraFields(ctx context.Context, action string, extraFields map[string]string) []zapcore.Field {
	zapFields := []zapcore.Field{}
	for key, value := range extraFields {
		zapFields = append(zapFields, zap.String(key, value))
	}

	requestID := appContext.GetRequestID(ctx)
	zapFields = append(zapFields, zap.String("request_id", requestID), zap.String("action", action))
	return zapFields
}

func getConfig(encoding string) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(getLogLevel()),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "timestamp",
			NameKey:    "logger",
			MessageKey: "message",
			CallerKey:  "caller",
			LevelKey:   "level",

			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     millisecondsTimeEnconder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func getLogLevel() zapcore.Level {
	logLevelConfig := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	if logLevelConfig == "DEBUG" {
		return zap.DebugLevel
	}

	if logLevelConfig == "INFO" {
		return zap.InfoLevel
	}

	if logLevelConfig == "WARN" {
		return zap.WarnLevel
	}

	if logLevelConfig == "ERROR" {
		return zap.ErrorLevel
	}

	return zap.InfoLevel
}

func millisecondsTimeEnconder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	nanos := t.UnixNano()
	millis := float64(nanos) / float64(time.Millisecond)
	enc.AppendInt64(int64(millis))
}
