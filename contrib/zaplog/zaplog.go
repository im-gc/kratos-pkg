package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(config *zap.Config)

// NewDefaultLogger ... 创建基于 zap 的日志实例, 并设置 zap 全局默认日志
func NewDefaultLogger(opts ...Option) *zap.Logger {
	// 优化日志的输出体验
	zapCfg := zap.NewProductionConfig()
	// 防止 kratos log 的 key 重复; took https://github.com/go-kratos/kratos/issues/1722
	zapCfg.EncoderConfig.TimeKey = ""
	zapCfg.EncoderConfig.MessageKey = ""
	zapCfg.EncoderConfig.CallerKey = ""
	// 默认等级为 debug
	// 如果有需求修改等级，可以用 WithLevel
	// 或者在 go-kratos 初始化 log 时，用 log.NewFilter() 进行过滤
	zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	for _, opt := range opts {
		opt(&zapCfg)
	}
	zapLog, _ := zapCfg.Build()
	zap.ReplaceGlobals(zapLog)
	return zapLog
}

// New ... alias for NewDefaultLogger
var New = NewDefaultLogger

// WithLevel ... 配置 zap-log 默认的日志等级
func WithLevel(level string) Option {
	return func(config *zap.Config) {
		l, err := zapcore.ParseLevel(level)
		if err == nil {
			config.Level.SetLevel(l)
		}
	}
}

func WithOutputPaths(outputPaths []string) Option {
	return func(config *zap.Config) {
		if len(outputPaths) > 0 {
			config.OutputPaths = outputPaths
		}
	}
}
