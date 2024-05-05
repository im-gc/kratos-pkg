package trace

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TracingConf interface {
	GetEndpoint() string
	GetCustomName() string
}

type tracingConfig struct {
	endpoint    string
	serviceName string
	fraction    float64
}

type Option func(tr *tracingConfig)

// InitTracer ... 初始化 tracer, 可以自定义 endpoint, service_name, ratio_based
func InitTracer(trConf TracingConf, opts ...Option) {
	var (
		tc = &tracingConfig{
			serviceName: trConf.GetCustomName(),
			endpoint:    trConf.GetEndpoint(),
			// rate based on the parent span to 100%
			fraction: 1.0,
		}
		host, _ = os.Hostname()
	)

	for _, opt := range opts {
		opt(tc)
	}

	if tc.endpoint == "" {
		return
	}
	if tc.serviceName == "" {
		tc.serviceName = host
	}

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(trConf.GetEndpoint())))
	if err != nil {
		return
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the `tracingConfig`
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(tc.fraction))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(tc.serviceName),
			attribute.String("host", host),
		)),
	)
	otel.SetTracerProvider(tp)
}

// WithEndpoint ... 自定义 jaeger 入口
func WithEndpoint(endpoint string) Option {
	return func(tr *tracingConfig) {
		tr.endpoint = endpoint
	}
}

// WithServiceName ...自定义上报服务名
func WithServiceName(name string) Option {
	return func(tr *tracingConfig) {
		tr.serviceName = name
	}
}

// WithRatioBased ... 设置采样比率
func WithRatioBased(fraction float64) Option {
	return func(tr *tracingConfig) {
		tr.fraction = fraction
	}
}
