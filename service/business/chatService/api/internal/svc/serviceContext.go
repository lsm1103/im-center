package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"im-center/service/business/chatService/api/internal/config"
	"im-center/service/model/cache"
	"im-center/common/nodeHandler"
)

type ServiceContext struct {
	Config config.Config
	RpcU *nodeHandler.RpcU
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RpcU: nodeHandler.NewRpcU(
			cache.NewRedisCache(redis.New(c.Redis.Host, func(r *redis.Redis) {
				r.Type = c.Redis.Type
				//r.Pass = c.Redis.Pass
			}) ) ),
	}
}
