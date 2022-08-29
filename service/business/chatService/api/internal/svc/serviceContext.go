package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"im-center/service/business/chatService/api/internal/config"
	"im-center/service/business/chatService/api/internal/utils"
	"im-center/service/model/cache"
)

type ServiceContext struct {
	Config config.Config
	RpcU *utils.RpcU
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RpcU: utils.NewRpcU(
			cache.NewRedisCache(redis.New(c.Redis.Host, func(r *redis.Redis) {
				r.Type = c.Redis.Type
				//r.Pass = c.Redis.Pass
			}) ) ),
	}
}
