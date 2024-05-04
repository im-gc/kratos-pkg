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

func Debug(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Debug(msg)
}

func Debugf(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Debugf(format, msg...)
}

func Info(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Info(msg)
}

func Infof(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Infof(format, msg...)
}

func Warn(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Warn(msg)
}

func Warnf(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Warnf(format, msg...)
}

func Error(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Error(msg)
}

func Errorf(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Errorf(format, msg...)
}

func Fatal(ctx context.Context, msg interface{}) {
	globalHelp.WithContext(ctx).Fatal(msg)
}

func Fatalf(ctx context.Context, format string, msg ...interface{}) {
	globalHelp.WithContext(ctx).Fatalf(format, msg...)
}
