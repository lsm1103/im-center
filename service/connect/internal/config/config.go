package config

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest"
	"im-center/common/ipTools"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	Redis struct {
		Host string
		Type string
		//Pass string
	}
	ServerIp string `json:",optional"`

	ConnRWSeparate bool `json:",default=false,optional"`
	IsDistributed bool `json:",default=false,optional"`
	ConnectInfoTimeout int `json:",default=300,optional"`	// 单位：秒, 默认5分钟
	ServerTimeout int `json:",default=300,optional"`	// 单位：秒, 默认5分钟
	RpcPort int `json:",default=2001"`	// rpc端口：默认2001
	TCPListenAddr string `json:",default=0.0.0.0:3001"`
}

func (c *Config) GetServerIp() string {
	if c.ServerIp == ""{
		c.ServerIp = fmt.Sprintf("%s:%d", ipTools.GetLocalIp(), c.RpcPort)
	}
	return c.ServerIp
}