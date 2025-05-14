package job_test

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/im-gc/kratos-pkg/contrib/job"
	"github.com/im-gc/kratos-pkg/contrib/trace"
	"github.com/im-gc/kratos-pkg/contrib/zaplog"
	"testing"
)

type testTracing struct {
	endpoint   string
	customName string
}

func (r testTracing) GetEndpoint() string {
	return r.endpoint
}

func (r testTracing) GetCustomName() string {
	return r.customName
}

func TestNewJobExecutor(t *testing.T) {
	logger := log.With(
		job.NewLogger(zaplog.NewDefaultLogger()),
		"timestamp", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	//log.SetLogger(logger)
	log := log.NewHelper(logger)

	trConf := &testTracing{
		endpoint:   "http://172.18.156.241:14268/api/traces",
		customName: "myapp",
	}

	trace.InitTracer(trConf)

	exec := job.NewJobExecutorWithTracer(logger, &job.Config{
		Enabled:      true,
		ServerAddr:   "http://172.18.156.47:8080/xxl-job-admin",
		AccessToken:  "qaz",
		ExecutorPort: 9999,
		RegistryKey:  "test-executor",
	}, nil)

	go func() {
		ctx := context.Background()
		exec.Start(ctx)
		log.WithContext(ctx).Infof("job executor started")
		defer exec.Stop(ctx)
	}()

	go func() {
		exec.RegisterTask("dym-task", func(ctx context.Context, payload string) (string, error) {
			log.WithContext(ctx).Infof("dym-task payload %s", payload)
			err := fmt.Errorf("一个错误")
			log.WithContext(ctx).Error(err)
			return "error", err
		})
	}()

	select {}
}
