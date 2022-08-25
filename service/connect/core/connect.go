/*
连接层的核心模块，包含多协议的长连接实现，连接的管理，关于连接、消息的最基础功能；
业务层可以在此基础上实现自己的业务，如聊天业务、客服业务、游戏业务、推送业务等；
*/
package core

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"im-center/common/serverInfo"
	"im-center/service/connect/core/connectManager"
	"im-center/service/connect/core/init/distributed"
	"im-center/service/connect/internal/config"
	"im-center/service/connect/internal/types"
	"im-center/service/connect/rpc/connectclient"
	"im-center/service/model/cache"
)

type ConnectServer struct {
	logx.Logger
	Config      *config.Config
	Distributed *distributed.Distributed
	cm *connectManager.ConnectManager
}

func NewConnectServer(cfg *config.Config, cache_ *cache.RedisCache) *ConnectServer {
	distributed_ := distributed.NewDistributed(cfg)
	if cfg.IsDistributed{
		distributed_.Start(cache_)
	}

	cm := connectManager.GetCM()
	cm.SetRelyOn(cfg, distributed_)
	return &ConnectServer{
		Logger: logx.WithContext(context.Background()),
		Config: cfg,
		Distributed: distributed_,
		cm: cm,
	}
}

func (c *ConnectServer) Start() {
	go c.cm.Start()
	//go tcp.TCPServer(&c.Config)
}

func (c *ConnectServer) Stop() {
	c.cm.Stop()
	//tcp.StopTCPServer()
}

func (c *ConnectServer) SendOneMsg_(req *types.SendOneMsgReq) error {
	//判断是否在本节点
	if c.cm.ExistConnect(&types.ConnectUid{ UserId: req.ReceiverId, DeviceId: req.ReceiverDeviceId}) {
		return c.cm.SendOneMsg(req)
	}

	//判断是否开启了多节点
	if !c.Distributed.IsDistributed() {
		return c.cm.SendOneMsg(req)
	}else {
		//如果开启了多节点，则通过rpc发送消息
		msgCenter, err := json.Marshal(req.MsgContent)
		if err != nil {
			c.Errorf("json.Marshal error: %s", err.Error())
		}
		_, err = c.Distributed.GetRpcNode(req.ReceiverDeviceId).SendOneMsg(
			context.Background(),
			&connectclient.SendOneMsgReq{
				SenderType:       req.SenderType,
				SenderId:         req.SenderId,
				SenderDeviceId:   req.SenderDeviceId,
				ReceiverId:       req.ReceiverId,
				ReceiverDeviceId: req.ReceiverDeviceId,
				ParentId:         req.ParentId,
				SendTime:         req.SendTime,
				MsgType:          req.MsgType,
				MsgContent: 	  string(msgCenter),
			})
		if err != nil {
			c.Errorf("send msg to user %s:%s error: %s", req.ReceiverId, req.ReceiverDeviceId, err.Error())
		}
	}
	return nil
}
func (c *ConnectServer) SendManyMsg_(req *types.SendManyMsgReq) error {
	//todo 可以做同节点消息聚合处理
	//拿到所有在本节点的消息，执行发送消息逻辑
	if !c.Distributed.IsDistributed(){
		return c.cm.SendManyMsg(req)
	}
	for _, userId := range req.ReceiverId {
		user, err := c.Distributed.Cache.GetConnectListByUser(userId)
		if err != nil {
			return err
		}
		for _, item := range user {
			err = c.SendOneMsg(&types.SendOneMsgReq{
				SenderType:       req.SenderType,
				SenderId:         req.SenderId,
				SenderDeviceId:   req.SenderDeviceId,
				ReceiverId:       item.UserId,
				ReceiverDeviceId: item.DeviceId,
				ParentId:         req.ParentId,
				SendTime:         req.SendTime,
				MsgType:          req.MsgType,
				MsgContent:       req.MsgContent,
			}, false)
			if err != nil {
				c.Errorf("send msg to user %s error: %s", item.UserId, err.Error())
			}
		}
	}

	return nil
}

//todo 所有外部调用，在rpc层面要记录消息、和业务数据
func (c *ConnectServer) SendOneMsg(req *types.SendOneMsgReq, isLocal bool) error {
	//判断是否在本节点
	if c.cm.ExistConnect(&types.ConnectUid{ UserId: req.ReceiverId, DeviceId: req.ReceiverDeviceId}) || !c.Distributed.IsDistributed() {
		return c.cm.SendOneMsg(req)
	} else if isLocal {
		return errors.New("not exist connect")
	}

	//开启了多节点
	list, err := c.Distributed.Cache.GetConnectListByUser(req.ReceiverId)
	if err != nil {
		return err
	}
	for _,info := range list {
		msgCenter, err := json.Marshal(req.MsgContent)
		if err != nil {
			c.Errorf("json.Marshal error: %s", err.Error())
		}
		//通过rpc发到相应的节点
		_, err = c.Distributed.GetRpcNode(info.ServerIp).SendOneMsg(
			context.Background(),
			&connectclient.SendOneMsgReq{
				SenderType:       req.SenderType,
				SenderId:         req.SenderId,
				SenderDeviceId:   req.SenderDeviceId,
				ReceiverId:       info.UserId,
				ReceiverDeviceId: info.DeviceId,
				ParentId:         req.ParentId,
				SendTime:         req.SendTime,
				MsgType:          req.MsgType,
				MsgContent: 	  string(msgCenter),
				IsLocal: 		  true,
			})
		if err != nil {
			//todo 某个连接发送错误的返回处理
			c.Errorf("send msg to user %s error: %s", info.UserId, err.Error())
		}
	}
	return nil
}

func (c *ConnectServer) SendManyMsg(req *types.SendManyMsgReq) error {
	//拿到所有在本节点的消息，执行发送消息逻辑
	if !c.Distributed.IsDistributed(){
		return c.cm.SendManyMsg(req)
	}
	//通过rpc把消息发到相应的节点，批量走SendManyMsg，单条走SendOneMsg
	for _,userId := range req.ReceiverId {
		list, err := c.Distributed.Cache.GetConnectListByUser(userId)
		if err != nil {
			c.Errorf("get user %s connect list by user error: %s", userId, err.Error())
		}
		for _,info := range list {
			msgCenter, err := json.Marshal(req.MsgContent)
			if err != nil {
				c.Errorf("json.Marshal error: %s", err.Error())
			}
			_, err = c.Distributed.GetRpcNode(info.ServerIp).SendOneMsg(
				context.Background(),
				&connectclient.SendOneMsgReq{
					SenderType:       req.SenderType,
					SenderId:         req.SenderId,
					SenderDeviceId:   req.SenderDeviceId,
					ReceiverId:       info.UserId,
					ReceiverDeviceId: info.DeviceId,
					ParentId:         req.ParentId,
					SendTime:         req.SendTime,
					MsgType:          req.MsgType,
					MsgContent: 	  string(msgCenter),
					IsLocal: 		  true,
				})
			if err != nil {
				c.Errorf("send msg to user %s error: %s", info.UserId, err.Error())
			}
		}

	}
	return nil
}

func (c *ConnectServer) OffConnect(req *types.ConnectUid) error {
	//判断是否在本节点
	if c.cm.ExistConnect(req) || !c.Distributed.IsDistributed() {
		return c.cm.OffConnect(req)
	}

	info, err := c.Distributed.Cache.GetConnectInfo(req.UserId, req.DeviceId)
	if err != nil {
		return err
	}
	_, err = c.Distributed.GetRpcNode(info.ServerIp).OffConnect(
		context.Background(),
		&connectclient.OffConnectReq{
			UserId: info.UserId,
			DeviceId: info.DeviceId,
		} )
	if err != nil {
		return err
	}
	return nil
}

func (c *ConnectServer) ExistConnect(req *types.ConnectUid) bool {
	//判断是否在本节点
	r := c.cm.ExistConnect(req)
	if !r && c.Distributed.IsDistributed(){
		return c.Distributed.Cache.ExistConnect(req.DeviceId, req.UserId)
	}
	return r
}

func (c *ConnectServer) GetConnectInfo(req *types.ConnectUid) (resp *types.ConnectInfoResp, err error) {
	if c.Distributed.IsDistributed(){
		info, err := c.Distributed.Cache.GetConnectInfo(req.DeviceId, req.UserId)
		if err != nil {
			return nil, err
		}
		return &types.ConnectInfoResp{
			UserId:         info.UserId,
			DeviceId:       info.DeviceId,
			ServerIp:       info.ServerIp,
			ConnectIp:      info.Addr,
			RegisterTime:   info.RegisterTime,
			HeartbeatTime:  info.HeartbeatTime,
			UnRegisterTime: info.UnRegisterTime,
			DeviceInfo:     info.DeviceInfo,
		}, nil
	}
	return c.cm.GetConnectInfo(req)
}

func (c *ConnectServer) GetConnectList(req *types.GetConnectListReq) (resp *types.ConnectListResp, err error) {
	if c.Distributed.IsDistributed(){
		list, err := c.Distributed.Cache.GetConnectList(req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}
		resp_ := []types.ConnectItem{}
		for _,item := range list {
			resp_ = append(resp_, types.ConnectItem{
				UserId:         item.UserId,
				DeviceId:       item.DeviceId,
				ServerIp:       item.ServerIp,
				ConnectIp:      item.Addr,
				RegisterTime:   item.RegisterTime,
				HeartbeatTime:  item.HeartbeatTime,
				UnRegisterTime: item.UnRegisterTime,
				DeviceInfo:     item.DeviceInfo,
			})
		}
		return &types.ConnectListResp{
			ConnectList: resp_,
		}, nil
	}
	return c.cm.GetConnectList(req)
}

func (c *ConnectServer) GetOnlineUserList(req *types.GetOnlineUserListReq) (resp *types.OnlineUserListResp, err error) {
	if ! c.Distributed.IsDistributed(){
		return c.cm.GetOnlineUserList(req)
	}
	// todo 分布式部分
	return nil, nil
}

func (c *ConnectServer) GetUserConnectList(req *types.GetUserConnectListReq) (resp *types.ConnectListResp, err error) {
	if c.Distributed.IsDistributed(){
		list, err := c.Distributed.Cache.GetConnectListByUser(req.UserId)
		if err != nil {
			return nil, err
		}
		resp_ := []types.ConnectItem{}
		for _,item := range list {
			resp_ = append(resp_, types.ConnectItem{
				UserId:         item.UserId,
				DeviceId:       item.DeviceId,
				ServerIp:       item.ServerIp,
				ConnectIp:      item.Addr,
				RegisterTime:   item.RegisterTime,
				HeartbeatTime:  item.HeartbeatTime,
				UnRegisterTime: item.UnRegisterTime,
				DeviceInfo:     item.DeviceInfo,
			})
		}
		return &types.ConnectListResp{
			ConnectList: resp_,
		}, nil
	}
	return c.cm.GetUserConnectList(req)
}

func (c *ConnectServer) GetServiceInfo(req *types.GetServerInfoReq) (resp *types.ServerInfoResp, err error) {
	if c.Distributed.IsDistributed(){
		return c.Distributed.GetServiceInfo()
	}
	serviceInfo := serverInfo.GetServiceInfo()
	sysInfo := serverInfo.GetSysInfo()
	return &types.ServerInfoResp{
		Server: []types.ServerItem{
			{
				ServerId:     c.Config.GetServerIp(),
				ServerInfo:   types.ServerInfo{
					NumGoroutine:     serviceInfo.NumGoroutine,
					AllocMemory:      serviceInfo.AllocMemory,
					TotalAllocMemory: serviceInfo.TotalAllocMemory,
					SysMemory:        serviceInfo.SysMemory,
					NumGC: 			  serviceInfo.NumGC,
				},
				BusinessInfo: *c.cm.GetBusinessInfo(),
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
			},
		},
	},nil
}
