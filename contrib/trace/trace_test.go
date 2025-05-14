package trace_test

import (
	"github.com/im-gc/kratos-pkg/contrib/trace"
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

func TestInitTrace(t *testing.T) {
	trConf := &testTracing{
		endpoint:   "http://localhost:14268/api/traces",
		customName: "myapp",
	}

	trace.InitTracer(trConf)
}
