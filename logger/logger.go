package logger

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	globalLogger log.Logger
	globalHelp   log.Helper
)

func GlobalLogger() log.Logger {
	return globalLogger
}

func With(kv ...interface{}) {
	globalLogger = log.With(globalLogger, kv...)
}

func Debug(msg interface{}) {
	globalHelp.Debug(msg)
}

func DebugWithContext(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Debug(msg)
}

func Debugf(format string, msg ...interface{}) {
	globalHelp.Debugf(format, msg...)
}

func DebugfWithContext(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Debugf(format, msg...)
}

func Info(msg interface{}) {
	globalHelp.Info(msg)
}

func InfoWithContext(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Info(msg)
}

func Infof(format string, msg ...interface{}) {
	globalHelp.Infof(format, msg...)
}

func InfofWithContext(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Infof(format, msg...)
}

func Warn(msg interface{}) {
	globalHelp.Warn(msg)
}

func WarnWithContext(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Warn(msg)
}

func Warnf(format string, msg ...interface{}) {
	globalHelp.Warnf(format, msg...)
}

func WarnfWithContext(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Warnf(format, msg...)
}

func Error(msg interface{}) {
	globalHelp.Error(msg)
}

func ErrorWithContext(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Error(msg)
}

func Errorf(format string, msg ...interface{}) {
	globalHelp.Errorf(format, msg...)
}

func ErrorfWithContext(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Errorf(format, msg...)
}

func Fatal(msg interface{}) {
	globalHelp.Fatal(msg)
}

func FatalWithContext(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Fatal(msg)
}

func Fatalf(format string, msg ...interface{}) {
	globalHelp.Fatalf(format, msg...)
}

func FatalfWithContext(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Fatalf(format, msg...)
}
