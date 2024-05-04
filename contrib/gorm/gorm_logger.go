package gorm

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	debug                 bool
	dbLog                 *log.Helper
	SlowThreshold         time.Duration
	SourceField           string
	SkipCallerLookup      bool
	SkipErrRecordNotFound bool
}

func NewGormLogger(dbLog *log.Helper, hasDebug bool) *Logger {
	return &Logger{
		debug:                 hasDebug,
		dbLog:                 dbLog,
		SlowThreshold:         500 * time.Millisecond, // 500毫秒查询 + 500毫秒业务响应 = 1s 用户最佳体验loading之内，超过则属于慢查询
		SkipCallerLookup:      false,
		SkipErrRecordNotFound: true,
	}
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.dbLog.WithContext(ctx).Info(s, args)
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.dbLog.WithContext(ctx).Warn(s, args)
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.dbLog.WithContext(ctx).Error(s, args)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.debug {
		return
	}
	clog := l.dbLog.WithContext(ctx)

	elapsed := time.Since(begin)
	timeUsed := float64(elapsed.Nanoseconds()) / 1e6

	fields := make([]interface{}, 0)
	fields = append(fields, "timeUsed(ms):", timeUsed)

	switch {
	case err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound):
		sql, rows := fc()
		fields = append(fields, "error:", err.Error())
		fields = append(fields, "sql:", sql)
		fields = append(fields, "rows:", rows)
		clog.Error(fields)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		fields = append(fields, "sql:", sql)
		fields = append(fields, "SLOW SQL:", l.SlowThreshold)
		fields = append(fields, "rows:", rows)
		clog.Warn(fields)
	default:
		sql, rows := fc()
		fields = append(fields, "sql:", sql)
		fields = append(fields, "rows:", rows)
		clog.Debug(fields)
	}
}
