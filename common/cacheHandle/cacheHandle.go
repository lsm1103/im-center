package cacheHandle

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
)
var (
	// ErrNotFound is an alias of predict.ErrNotFound.
	ErrNotFound = errors.New("predict: no rows in result set")

	// can't use one SingleFlight per conn, because multiple conns may share the same cache key.
	singleFlights = syncx.NewSingleFlight()
	stats         = cache.NewStat("predict")
)

type  CecheHandle struct {
	cc cache.Cache
}

func NewCecheHandle(clusterConf cache.ClusterConf) *CecheHandle {
	return &CecheHandle{
		cc : cache.New(clusterConf, singleFlights, stats, ErrNotFound),
	}
}

//读：先读缓存，如果找不到再读真实存储，然后更新缓存；
func (c CecheHandle)Get(val interface{}, key string, query func(val interface{}) error) error {
	return c.cc.Take(val, key, query)
	//return c.cc.TakeWithExpire(val, key, query), expire time.Duration
}

//写：先写入真实存储，再删除缓存（考虑脏数据可以延迟双删）；
func (c CecheHandle) Exec(exec func() error, keys ...string) (err error) {
	if err = exec(); err != nil {
		return err
	}
	if err = c.cc.Del(keys...); err != nil {
		return err
	}
	return nil
}
