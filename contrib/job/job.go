package job

import (
	"context"
	"fmt"
	"github.com/go-basic/ipv4"
	"github.com/go-kratos/kratos/v2/log"
	http "github.com/go-kratos/kratos/v2/transport/http"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

var moduleName = "job/executor"

type JobExecutor struct {
	log        *log.Helper
	httpServer *http.Server
	exec       xxl.Executor
	tracer     trace.Tracer
}

func (r *JobExecutor) Start(ctx context.Context) error {
	if r.httpServer != nil {
		return nil
	}
	r.log.WithContext(ctx).Info("xxl job executor is starting")
	err := r.exec.Run()
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return err
	}
	r.log.WithContext(ctx).Warn("xxl job executor is started")
	return nil
}
func (r *JobExecutor) Stop(ctx context.Context) error {
	r.log.WithContext(ctx).Info("xxl job executor is stoping")
	r.exec.Stop()
	r.log.WithContext(ctx).Info("xxl job executor is stopped")
	return nil
}

// JobFunc (ctx, JSONString)
type JobFunc func(ctx context.Context, payload string) (string, error)

func (r *JobExecutor) RegisterTask(pattern string, fn JobFunc) {
	// 构建任务
	task := func(ctx context.Context, req *xxl.RunReq) string {
		var err error
		// 支持链路上报
		if r.tracer != nil {
			var (
				span trace.Span
				tn   = time.Now()
			)
			ctx, span = r.tracer.Start(ctx, fmt.Sprintf("%s %s", moduleName, pattern),
				trace.WithSpanKind(trace.SpanKindInternal))
			defer func() {
				if err != nil {
					span.RecordError(err, trace.WithTimestamp(tn))
				}
				setSpanAttrs(ctx, span, req)
				span.End()
			}()
		}
		ret, err := fn(ctx, req.ExecutorParams)
		if err != nil {
			r.log.WithContext(ctx).Error(err)
			return err.Error()
		}
		return ret
	}
	r.exec.RegTask(pattern, task)
	return
}

func registerHTTPServer(s *http.Server, exec xxl.Executor) {
	router := s.Route("/")

	router.POST("run", func(ctx http.Context) error {
		exec.RunTask(ctx.Response(), ctx.Request())
		return nil
	})
	router.POST("kill", func(ctx http.Context) error {
		exec.KillTask(ctx.Response(), ctx.Request())
		return nil
	})
	router.POST("log", func(ctx http.Context) error {
		exec.TaskLog(ctx.Response(), ctx.Request())
		return nil
	})
	router.POST("beat", func(ctx http.Context) error {
		exec.Beat(ctx.Response(), ctx.Request())
		return nil
	})
	router.POST("idleBeat", func(ctx http.Context) error {
		exec.IdleBeat(ctx.Response(), ctx.Request())
		return nil
	})
}

func NewJobExecutorWithTracer(logger log.Logger, conf JobConfig, httpServer *http.Server) *JobExecutor {
	jobExecutor := NewJobExecutor(logger, conf, httpServer)
	jobExecutor.tracer = otel.Tracer(moduleName)
	return jobExecutor
}

func NewJobExecutor(logger log.Logger, conf JobConfig, httpServer *http.Server) *JobExecutor {
	clog := log.NewHelper(
		log.With(
			logger, // 	log.NewFilter(logger, log.FilterLevel(log.LevelInfo)),
			"module", moduleName,
		),
	)
	//log.NewFilter(logger, log.FilterLevel(log.LevelInfo))

	executorIp := conf.GetExecutorIp()
	if executorIp == "" {
		executorIp = ipv4.LocalIP()
	}
	exec := xxl.NewExecutor(
		xxl.ServerAddr(conf.GetServerAddr()),                                   // 调度中心地址
		xxl.AccessToken(conf.GetAccessToken()),                                 // 请求令牌(默认为空)
		xxl.ExecutorIp(executorIp),                                             // 监听ip
		xxl.ExecutorPort(strconv.FormatInt(int64(conf.GetExecutorPort()), 10)), // 本地(执行器)端口 （如果端口不配置，则复用当前HTTP端口）
		xxl.RegistryKey(conf.GetRegistryKey()),                                 // 执行器名称
		xxl.SetLogger(&xxlLogger{clog}),
	)

	exec.Init()
	//设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: xxl.SuccessCode, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "这个是自定义日志handler",
			IsEnd:       true,
		}}
	})

	exec.RegTask("test", func(ctx context.Context, payload *xxl.RunReq) string {
		clog.WithContext(ctx).Infof("test task payload: %s", payload.ExecutorParams)
		return "success"
	})

	if httpServer != nil {
		registerHTTPServer(httpServer, exec)
	}

	return &JobExecutor{
		exec:       exec,
		log:        clog,
		httpServer: httpServer,
	}
}

func setSpanAttrs(ctx context.Context, span trace.Span, req *xxl.RunReq) {
	attrs := make([]attribute.KeyValue, 0)
	attrs = append(attrs, attribute.Key("jobId").Int64(req.JobID))
	attrs = append(attrs, attribute.Key("executorHandler").String(req.ExecutorHandler))
	attrs = append(attrs, attribute.Key("executorParams").String(req.ExecutorParams))
	attrs = append(attrs, attribute.Key("executorBlockStrategy").String(req.ExecutorBlockStrategy))
	attrs = append(attrs, attribute.Key("logId").Int64(req.LogID))
	attrs = append(attrs, attribute.Key("logDateTime").Int64(req.LogDateTime))
	attrs = append(attrs, attribute.Key("glueType").String(req.GlueType))
	attrs = append(attrs, attribute.Key("glueSource").String(req.GlueSource))
	attrs = append(attrs, attribute.Key("glueUpdatetime").Int64(req.GlueUpdatetime))
	attrs = append(attrs, attribute.Key("broadcastIndex").Int64(req.BroadcastIndex))
	attrs = append(attrs, attribute.Key("broadcastTotal").Int64(req.BroadcastTotal))
	span.SetAttributes(attrs...)
}

type xxlLogger struct {
	clog *log.Helper
}

func (l *xxlLogger) Info(format string, a ...interface{}) {
	l.clog.Infof(format, a...)
}
func (l *xxlLogger) Error(format string, a ...interface{}) {
	l.clog.Errorf(format, a...)
}
