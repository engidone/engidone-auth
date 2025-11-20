package log

import (
	"log"
	"time"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	baseLogger, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		log.Fatal(err)
	}
	logger := baseLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	sugar = logger
	defer logger.Sync()
}

func Info(args ...any) {
	sugar.Info(args...)
}

func Error(args ...any) {
	sugar.Error(args...)
}

func Warn(args ...any) {
	sugar.Warn(args...)
}

func Debug(args ...any) {
	sugar.Debug(args...)
}

func Infof(template string, args ...any) {
	sugar.Infof(template, args...)
}

func Fatal(args ...any) {
	sugar.Fatal(args...)
}

func Fatalf(template string, args ...any) {
	sugar.Fatalf(template, args...)
}

func Errorf(template string, args ...any) {
	sugar.Errorf(template, args...)
}

func Warnf(template string, args ...any) {
	sugar.Warnf(template, args...)
}

func Debugf(template string, args ...any) {
	sugar.Debugf(template, args...)
}

func Infow(msg string, keysAndValues ...any) {
	sugar.Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...any) {
	sugar.Errorw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...any) {
	sugar.Warnw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...any) {
	sugar.Debugw(msg, keysAndValues...)
}

type LogField = zap.Field

func String(key, value string) LogField {
	return zap.String(key, value)
}

func Int(key string, value int) LogField {
	return zap.Int(key, value)
}

func Int32(key string, value int32) LogField {
	return zap.Int32(key, value)
}

func Bool(key string, value bool) LogField {
	return zap.Bool(key, value)
}

func Err(err error) LogField {
	return zap.Error(err)
}

func Any(key string, value any) LogField {
	return zap.Any(key, value)
}

func Float64(key string, value float64) LogField {
	return zap.Float64(key, value)
}

func Int64(key string, value int64) LogField {
	return zap.Int64(key, value)
}

func Uint(key string, value uint) LogField {
	return zap.Uint(key, value)
}

func Uint32(key string, value uint32) LogField {
	return zap.Uint32(key, value)
}

func Uint64(key string, value uint64) LogField {
	return zap.Uint64(key, value)
}

func Duration(key string, value time.Duration) LogField {
	return zap.Duration(key, value)
}

func Time(key string, value time.Time) LogField {
	return zap.Time(key, value)
}
