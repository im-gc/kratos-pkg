package logger

import (
	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(GlobalLogger)

func GlobalLoad(level, serviceID, serviceName, serviceVersion string) error {
	// 优化日志的输出体验
	zapCfg := zap.NewProductionConfig()
	// 防止 kratos log 的 key 重复; took https://github.com/go-kratos/kratos/issues/1722
	zapCfg.EncoderConfig.TimeKey = ""
	zapCfg.EncoderConfig.MessageKey = ""
	zapCfg.EncoderConfig.CallerKey = ""
	zapLog, _ := zapCfg.Build()
	zap.ReplaceGlobals(zapLog)

	globalLogger = log.NewFilter(
		log.With(
			kratoszap.NewLogger(zapLog),
			"ts", log.DefaultTimestamp,
			"caller", log.Caller(6),
			"service.id", serviceID,
			"service.name", serviceName,
			"service.version", serviceVersion,
			"trace.id", tracing.TraceID(),
			"span.id", tracing.SpanID()),
		log.FilterLevel(log.ParseLevel(level)),
	)
	globalHelp = *log.NewHelper(globalLogger)
	return nil
}
