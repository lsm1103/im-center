package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"

	"im-center/service/connect/core"
	"im-center/service/connect/internal/config"
	"im-center/service/model/cache"
)

type ServiceContext struct {
	Config config.Config
	Cache *cache.RedisCache

	Cs *core.ConnectServer
}

func NewServiceContext(c config.Config) *ServiceContext {
	cache_ := cache.NewRedisCache(redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = c.Redis.Type
		//r.Pass = c.Redis.Pass
	}) )
	return &ServiceContext{
		Config: c,
		Cache: cache_,
	}
}
