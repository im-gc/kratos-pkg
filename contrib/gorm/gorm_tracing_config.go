package gorm

import oteltrace "go.opentelemetry.io/otel/trace"

type config struct {
	dbName         string
	tracerProvider oteltrace.TracerProvider
	alwaysOmitVars bool
}

// GormTracingOption is used to configure the client.
type GormTracingOption interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) GormTracingOption {
	return optionFunc(func(cfg *config) {
		cfg.tracerProvider = provider
	})
}

// WithDBName specified the database name to be used in span names
// since its not possible to extract this information from gorm
func WithDBName(name string) GormTracingOption {
	return optionFunc(func(cfg *config) {
		cfg.dbName = name
	})
}

// WithAlwaysOmitVariables makes the plugin always omit variable values from traces.
func WithAlwaysOmitVariables() GormTracingOption {
	return optionFunc(func(cfg *config) {
		cfg.alwaysOmitVars = true
	})
}
