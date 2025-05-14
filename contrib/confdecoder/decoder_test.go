package confdecoder_test

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/im-gc/kratos-pkg/contrib/confdecoder"
	"testing"
)

func TestNormalDecoder(t *testing.T) {

	t.Run("test linux .swp file", func(tt *testing.T) {
		if err := confdecoder.NormalDecoder(
			&config.KeyValue{Key: ".config.yaml.swp", Format: "swp"},
			map[string]interface{}{},
		); err != nil {
			tt.Fatal(err)
		}
	})

	t.Run("test Jetbrains IDE temp file", func(tt *testing.T) {
		if err := confdecoder.NormalDecoder(
			&config.KeyValue{Key: "config.yaml~", Format: "yaml~"},
			map[string]interface{}{},
		); err != nil {
			tt.Fatal(err)
		}
	})
}
