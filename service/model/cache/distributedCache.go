// 分布式模式下的缓存集群方法类
package cache

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"im-center/common/globalkey"
)

type (
	RedisCache struct {
		logx.Logger
		RdsC *redis.Redis
	}
	ConnectInfo struct {
		UserId         string `json:"user_id"`          // 用户Id
		DeviceId       string `json:"device_id"`        // 登录的平台Id app/web/ios
		Addr           string `json:"addr"`             // 客户端地址
		ServerIp       string `json:"server_ip"`        // 服务器ip
		RegisterTime   uint64 `json:"register_time"`    // 用户上次上线时间
		HeartbeatTime  uint64 `json:"heartbeat_time"`   // 用户上次心跳时间
		UnRegisterTime uint64 `json:"un_register_time"` // 用户上次下线时间
		DeviceInfo     string `json:"device_info"`      // 用户设备信息
	}
)

func NewRedisCache(rdsC *redis.Redis) *RedisCache {
	return &RedisCache{
		Logger: logx.WithContext(context.Background() ),
		RdsC: rdsC,
	}
}

func (d *RedisCache) NodeRegister(nodeId, data string, timeOut int) error {
	err := d.RdsC.Setex(globalkey.BuildKey(globalkey.CacheNodeIdKey, nodeId), data, timeOut)
	if err != nil {
		return err
	}
	return nil
}

func (d *RedisCache) GetNodeList() []string {
	r, err := d.RdsC.Keys(globalkey.CacheNodeIdMATCHKey)
	if err != nil {
		d.Logger.Errorf("获取节点keys失败：%v", err)
		return nil
	}
	nodes,err := d.RdsC.Mget(r...)
	if err != nil {
		d.Logger.Errorf("获取节点失败：%v", err)
		return nil
	}
	return nodes
}

func (d *RedisCache) SaveConnect(data *ConnectInfo, timeOut int) error {
	info, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = d.RdsC.SetnxEx(globalkey.BuildKey(globalkey.CacheConnectIdKey, data.UserId, data.DeviceId), string(info), timeOut)
	return err
}

func (d *RedisCache) DeleteConnect(userId string, deviceId string) error {
	_, err := d.RdsC.Del(globalkey.BuildKey(globalkey.CacheConnectIdKey, userId, deviceId))
	if err != nil {
		return err
	}
	return nil
}

func (d *RedisCache) ExistConnect(userId string, deviceId string) bool {
	exists, err := d.RdsC.Exists(globalkey.BuildKey(globalkey.CacheConnectIdKey, userId, deviceId))
	if err != nil {
		d.Logger.Errorf("判断连接是否存在失败：%v", err)
	}
	return exists
}

func (d *RedisCache) GetConnectInfo(userId string, deviceId string) (*ConnectInfo, error) {
	info, err := d.RdsC.Get(globalkey.BuildKey(globalkey.CacheConnectIdKey, userId, deviceId))
	if err != nil {
		return nil, err
	}
	if info == "" {
		return nil, nil
	}
	connInfo := &ConnectInfo{}
	err = json.Unmarshal([]byte(info), connInfo)
	if err != nil {
		d.Logger.Errorf("获取连接信息失败：%v", err)
		return nil, err
	}

	return connInfo, nil
}

func (d *RedisCache) GetConnectList(offset, limit uint64) ([]*ConnectInfo, error) {
	r, _, err := d.RdsC.Scan(offset, globalkey.CacheConnectIdMATCHKey, int64(limit))
	if err != nil {
		return nil, err
	}
	resp := []*ConnectInfo{}
	for _,item := range r {
		info := ConnectInfo{}
		err = json.Unmarshal([]byte(item), &info)
		if err != nil {
			return nil, err
		}
		resp = append(resp, &info)
	}
	return resp, nil
}

func (d *RedisCache) GetConnectListByUser(userId string) ([]*ConnectInfo, error) {
	r, err := d.RdsC.Keys(globalkey.BuildKey(globalkey.CacheConnectIdUserIdMATCHKey, userId) )
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		return nil, errors.New("没有找到该连接")
	}
	connectList, err := d.RdsC.Mget(r...)
	if err != nil {
		return nil, err
	}
	resp := []*ConnectInfo{}
	for _,item := range connectList {
		info := &ConnectInfo{}
		err = json.Unmarshal([]byte(item), info)
		if err != nil {
			return nil, err
		}
		resp = append(resp, info)
	}
	return resp, nil
}

func (d *RedisCache) ExpireConnect(userId string, deviceId string, seconds int) error {
	return d.RdsC.Expire(globalkey.BuildKey(globalkey.CacheConnectIdKey, userId, deviceId), seconds)
}