package contrib

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type DataConf interface {
	GetDriver() string
	GetSource() string
	// GetDebug ... alias IsDebug()
	GetDebug() bool
}

// NewGorm 创建 gorm 连接
// Deprecated: please use bt.baishancloud.com/resdev/contrib/gorm.NewGorm()
func NewGorm(c DataConf, ctxlog *log.Helper) (*gorm.DB, error) {
	var (
		driver gorm.Dialector
		opts   = &gorm.Config{
			// 默认不输出 log
			Logger: gormLogger.Discard,
		}
	)

	// only support MySQL.
	driver = mysql.New(mysql.Config{
		DSN:                      c.GetSource(),
		DisableDatetimePrecision: true,
		DontSupportRenameIndex:   true,
	})
	ctxlog.Debug("use gorm data-source from 'data.database.source'")

	// 覆盖 gorm logger，可以打印出 slow sql
	if c.GetDebug() {
		opts.Logger = NewGormLogger(ctxlog, c.GetDebug())
	}

	return gorm.Open(driver, opts)
}

type GormLogger struct {
	debug                 bool
	dbLog                 *log.Helper
	SlowThreshold         time.Duration
	SourceField           string
	SkipCallerLookup      bool
	SkipErrRecordNotFound bool
}

func NewGormLogger(dbLog *log.Helper, hasDebug bool) *GormLogger {
	return &GormLogger{
		debug:                 hasDebug,
		dbLog:                 dbLog,
		SlowThreshold:         500 * time.Millisecond, // 500毫秒查询 + 500毫秒业务响应 = 1s 用户最佳体验loading之内，超过则属于慢查询
		SkipCallerLookup:      false,
		SkipErrRecordNotFound: true,
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.dbLog.Info(s, args)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.dbLog.Warn(s, args)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.dbLog.Error(s, args)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.debug {
		return
	}
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
		l.dbLog.Error(fields)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		fields = append(fields, "sql:", sql)
		fields = append(fields, "SLOW SQL:", l.SlowThreshold)
		fields = append(fields, "rows:", rows)
		l.dbLog.Warn(fields)
	default:
		sql, rows := fc()
		fields = append(fields, "sql:", sql)
		fields = append(fields, "rows:", rows)
		l.dbLog.Debug(fields)
	}
}
