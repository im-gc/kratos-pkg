package gorm

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type DataConf interface {
	GetDriver() string
	GetSource() string
	// GetDebug ... alias IsDebug()
	GetDebug() bool
}

// Config 自定义 Gorm 配置项
type Config struct {
	driver   gorm.Dialector
	opts     *gorm.Config
	log      gormLogger.Interface
	tracing  bool
	hasDebug bool
}

type Option func(*Config)

// New 创建 gorm 连接
func New(c DataConf, ctxlog *log.Helper) (*gorm.DB, error) {
	var (
		driver gorm.Dialector
		opts   = &gorm.Config{
			Logger: gormLogger.Discard, // 默认不输出 log
		}
	)

	// 仅支持 MySQL 驱动
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

func NewWithOptions(c DataConf, opts ...Option) (*gorm.DB, error) {
	var o = &Config{
		hasDebug: c.GetDebug(),
		opts:     nil,
	}

	for _, opt := range opts {
		opt(o)
	}

	// 如果没有配置 driver 则创建默认的 mysql driver
	if o.driver == nil {
		o.driver = mysql.New(mysql.Config{
			DSN:                      c.GetSource(),
			DisableDatetimePrecision: true,
			DontSupportRenameIndex:   true,
		})
	}
	// 如果有配置 logger
	if o.log == nil {
		o.log = gormLogger.Discard
	}

	// 创建 gorm-config 并绑定参数
	o.opts = &gorm.Config{
		Logger:               o.log, // 默认关闭Logger
		CreateBatchSize:      1000,  // 默认 1000
		AllowGlobalUpdate:    false, // 默认不允许全表更新
		DisableAutomaticPing: false, // 默认不禁用自动ping (数据库连接保活)
	}

	db, err := gorm.Open(o.driver, o.opts)
	if err != nil {
		return nil, err
	}

	// 如果有配置 tracing 则启用 otel tracing
	if o.tracing {
		_ = db.Use(NewPlugin())
	}

	return db, nil
}

// WithDriver set gorm-driver.
func WithDriver(driver gorm.Dialector) Option {
	return func(c *Config) {
		c.driver = driver
	}
}

// WithLogger set gorm-logger and has debug logger writer..
func WithLogger(log *log.Helper, hasDebug bool) Option {
	return func(c *Config) {
		c.log = NewGormLogger(log, hasDebug)
	}
}

// WithTracing set gorm-tracing. used for opentracing.
func WithTracing() Option {
	return func(c *Config) {
		c.tracing = true
	}
}
