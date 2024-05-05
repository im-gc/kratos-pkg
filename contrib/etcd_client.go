package contrib

import clientv3 "go.etcd.io/etcd/client/v3"

// NewEtcdClient ... 根据 endpoints 创建 etcd client(v3)
func NewEtcdClient(endpoints []string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
}
