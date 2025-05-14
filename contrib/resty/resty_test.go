package resty_test

import (
	"github.com/im-gc/kratos-pkg/contrib/resty"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := resty.New(
		resty.WithDebug(true),
		resty.WithTimeout(5*time.Second),
	)

	resp, err := client.R().Get("http://localhost:2022/api/proxy/backend9091/baseinfo")

	t.Log(resp)
	t.Log(err)
}
