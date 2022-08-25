// 连接层分布式处理模块
package distributed

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"

	"im-center/common/globalkey"
	"im-center/common/serverInfo"
	"im-center/service/connect/internal/config"
	"im-center/service/connect/internal/types"
	"im-center/service/connect/rpc/connectclient"
	"im-center/service/model/cache"
)

type Distributed struct {
	logx.Logger
	config *config.Config
	Cache *cache.RedisCache
}

func NewDistributed(c *config.Config) *Distributed {
	return &Distributed{
		Logger: logx.WithContext(context.Background()),
		config: c,
	}
}

// 启动本模块
func (c *Distributed) Start(cache *cache.RedisCache) {
	c.Cache = cache
	c.NodeRegister(nil)
}

// 关闭本模块，关闭事件处理管道，
func (c *Distributed) Stop() {
}

//是否为本地节点
//func (c *Distributed) IsLocalNode(serverIp string) bool {
//	if serverIp == c.config.GetServerIp() {
//		return true
//	}
//	return false
//}
func (c *Distributed) IsDistributed() (resp bool) {
	//todo 实现分布式模块是否开启的有效期处理
	if c.Cache != nil {
		resp = true
	}
	//if c.Cache.GetNodeList() != nil {
	//	resp = true
	//}
	//c.Infof("IsDistributed:%+v", resp)
	return resp
}

func (c *Distributed) NodeRegister(info *types.ServerItem) {
	if info == nil {
		serviceInfo := serverInfo.GetServiceInfo()
		sysInfo := serverInfo.GetSysInfo()
		info = &types.ServerItem{
			ServerId:     c.config.GetServerIp(),
			ServerInfo:   types.ServerInfo{
				NumGoroutine:     serviceInfo.NumGoroutine,
				AllocMemory:      serviceInfo.AllocMemory,
				TotalAllocMemory: serviceInfo.TotalAllocMemory,
				SysMemory:        serviceInfo.SysMemory,
				NumGC: 			  serviceInfo.NumGC,
			},
			SysInfo:      types.SysInfo{
				NumCpu:          sysInfo.NumCpu,
				NumCpuUsage:     sysInfo.NumCpuUsage,
				NumRam:          sysInfo.NumRam,
				NumRamUsage:     sysInfo.NumRamUsage,
				NumDisk:         sysInfo.NumDisk,
				NumDiskUsage:    sysInfo.NumDiskUsage,
				NumNetwork:      sysInfo.NumNetwork,
				NumNetworkUsage: sysInfo.NumNetworkUsage,
			},
		}
	}

	data, err := json.Marshal(info)
	if err != nil {
		c.Errorf("json服务信息失败：%+v", err)
		return
	}
	err = c.Cache.NodeRegister(c.config.GetServerIp(), string(data), c.config.ServerTimeout)
	if err != nil {
		c.Errorf("节点注册失败：%+v", err)
		return
	}
}

func (c *Distributed) GetServiceInfo() (resp *types.ServerInfoResp, err error) {
	nodeList := c.Cache.GetNodeList()
	var resp_ []types.ServerItem
	for _,item := range nodeList {
		tmp := &types.ServerItem{}
		err = json.Unmarshal([]byte(item), tmp)
		if err != nil {
			c.Errorf("获取分布式节点列表：%+v,%+v", item, err)
			return nil, err
		}
		resp_ = append(resp_, *tmp)
	}
	c.Errorf("获取分布式节点列表：%+v", resp_)
	return &types.ServerInfoResp{
		Server: resp_,
	}, nil
}

// 获取分布式节点
func (c *Distributed) GetRpcNode(rpcIp string) connectclient.Connect {
	return connectclient.NewConnect(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{rpcIp},
		NonBlock: true,
	}))
}

func (c *Distributed) EventHandler(event *types.DistributedEvent) {
	if !c.IsDistributed() {
		return
	}
	//处理事件
	switch event.EventType {
	case globalkey.EventTypeOnline:
		err := c.onlineEventLogic(event.ConnectData)
		if err != nil {
			c.Logger.Errorf("上线事件处理逻辑失败：%+v", err)
			return
		}
	case globalkey.EventTypeOffline:
		err := c.offlineEventLogic(event.ConnectData)
		if err != nil {
			c.Logger.Errorf("下线事件处理逻辑失败：%+v", err)
			return
		}
	case globalkey.EventTypeRenewal:
		err := c.renewalEventLogic(event.ConnectData)
		if err != nil {
			c.Logger.Errorf("连接续期事件处理逻辑失败：%+v", err)
			return
		}
	}
}

// 连接层分布式处理模块上线事件处理逻辑
func (c *Distributed) onlineEventLogic(data *types.ConnectData) error {
	//把连接写到公共缓存集群中
	err := c.Cache.SaveConnect(&cache.ConnectInfo{
		Addr:           data.Addr,
		DeviceId:       data.DeviceId,
		UserId:         data.UserId,
		RegisterTime:   data.RegisterTime,
		HeartbeatTime:  data.HeartbeatTime,
		ServerIp:       data.ServerIp,
	}, c.config.ConnectInfoTimeout)
	if err != nil {
		return err
	}
	return nil
}

// 连接层分布式处理模块下线事件处理逻辑
func (c *Distributed) offlineEventLogic(connect *types.ConnectData) error {
	//把连接从公共缓存集群中删除
	return c.Cache.DeleteConnect(connect.UserId, connect.DeviceId)
}

// 连接层分布式处理模块连接续期事件处理逻辑
func (c *Distributed) renewalEventLogic(connect *types.ConnectData) error {
	//把连接从公共缓存集群中续期
	return c.Cache.ExpireConnect(connect.UserId, connect.DeviceId, c.config.ConnectInfoTimeout)
}