package job

import (
	"google.golang.org/protobuf/types/known/durationpb"
)

type JobConfig interface {
	GetEnabled() bool
	GetServerAddr() string
	GetAccessToken() string
	GetTimeout() *durationpb.Duration
	GetExecutorIp() string
	GetExecutorPort() int32
	GetRegistryKey() string
}

type Config struct {
	Enabled      bool
	ServerAddr   string
	AccessToken  string
	Timeout      *durationpb.Duration
	ExecutorIp   string
	ExecutorPort int32
	RegistryKey  string
}

func (c *Config) GetEnabled() bool {
	return c.Enabled
}
func (c *Config) GetServerAddr() string {
	return c.ServerAddr
}
func (c *Config) GetAccessToken() string {
	return c.AccessToken
}
func (c *Config) GetTimeout() *durationpb.Duration {
	return c.Timeout
}
func (c *Config) GetExecutorIp() string {
	return c.ExecutorIp
}
func (c *Config) GetExecutorPort() int32 {
	return c.ExecutorPort
}
func (c *Config) GetRegistryKey() string {
	return c.RegistryKey
}
