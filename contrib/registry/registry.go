package registry

import (
	"errors"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// RegistryConf ... 默认的服务发现和注册中心配置项
type RegistryConf interface {
	GetEndpoints() []string
	GetEnabled() bool
}

var ErrEndpointNotExist = errors.New("ERR_ENDPOINT_NOT_EXIST")

// NewEtcdClient .... init etcd client/v3
func NewEtcdClient(conf RegistryConf) (*clientv3.Client, func(), error) {
	if !conf.GetEnabled() {
		return nil, func() {}, nil
	}
	// check endpoints
	if conf.GetEndpoints() == nil || len(conf.GetEndpoints()) <= 0 {
		return nil, func() {}, ErrEndpointNotExist
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: conf.GetEndpoints(),
	})

	cleanup := func() {
		if cli != nil {
			_ = cli.Close()
		}
	}

	if err != nil {
		return nil, cleanup, err
	}

	return cli, cleanup, err
}

// NewRegistrar ... init ectd Registry
func NewRegistrar(client *clientv3.Client) registry.Registrar {
	if client == nil {
		return nil
	}
	return etcd.New(client)
}

// NewDiscovery ... init etcd Discovery
func NewDiscovery(client *clientv3.Client) registry.Discovery {
	if client == nil {
		return nil
	}
	return etcd.New(client)
}
