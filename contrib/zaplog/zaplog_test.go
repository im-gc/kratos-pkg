package zaplog

import "testing"

func TestOption(t *testing.T) {
	logger := NewDefaultLogger(
		WithOutputPaths([]string{"/tmp/logs"}),
		WithLevel("debug"),
	)
	logger.Info("info logs")
	logger.Error("error logs")
}
