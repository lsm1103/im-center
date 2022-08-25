package ws

import (
	"context"
	"im-center/service/connect/internal/config"
	"runtime/debug"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"

	"im-center/common/globalkey"
	"im-center/common/tool"
	"im-center/service/connect/core/connect"
	"im-center/service/connect/core/connectManager"
	"im-center/service/connect/internal/types"
)

type wsConnect struct {
	logx.Logger
	ctx    context.Context
	config      *config.Config

	Socket        *websocket.Conn // 连接
	SocketLock    sync.RWMutex    // 读写锁
	Send          chan []byte     // 待发送的数据
	UserId        string          // 用户Id
	DeviceId      string          // 登录的平台Id app/web/ios
	Addr          string          // 客户端地址
	RegisterTime  uint64          // 用户上次上线时间
	HeartbeatTime uint64          // 用户上次心跳时间
	SetCacheTime  uint64          // 用户上次设置缓存时间
	ServerIp	  string          // 服务器ip
	//UnRegisterTime uint64       // 用户上次下线时间
	//DeviceInfo    string		  // 用户设备信息
}

var cm = connectManager.GetCM()

func NewWsConnect(ctx context.Context, config *config.Config, in types.WsConnect) connect.Connect {
	return &wsConnect{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		config: config,

		Socket:        in.Socket,
		Send:          make(chan []byte, 100),
		Addr:          in.Addr,
		UserId:        in.UserId,
		DeviceId:      in.DeviceId,
		RegisterTime:  in.CurrentTime,
		HeartbeatTime: in.CurrentTime,
		SetCacheTime:  in.CurrentTime,
		ServerIp: 	   config.GetServerIp(),
	}
}

func (c *wsConnect) Read() {
	defer func() {
		c.Error("defer 关闭Read协程", c.UserId, c.DeviceId)
		if r := recover(); r != nil {
			c.Error("write stop", string(debug.Stack()), r)
		}
		if c.Send != nil {
			// 连接下线处理
			cm.Unregister <- c
		}
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			c.Error("读取客户端数据错误", c.Addr, err)
			return
		}
		// 处理程序
		connect.ProcessData(c, message)
	}
}

func (c *wsConnect) Write() {
	defer func() {
		c.Errorf("defer, 关闭Write协程 %+v", c)
		if r := recover(); r != nil {
			c.Errorf("write stop", string(debug.Stack()), r)
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				c.Errorf("Client发送数据 关闭连接:%+s, ok:%+v", c.Addr, ok)
				return
			}
			c.socketWriteMsg(message)
		}
	}
}

func (c *wsConnect) ActiveClose() {
	if c.Socket != nil {
		err := c.Socket.Close()
		if err != nil {
			c.Errorf("关闭连接错误", c.Addr, err)
		}
	}
}

func (c *wsConnect) PassiveClose() {
	if c.Send != nil {
		close(c.Send)
	}
}

func (c *wsConnect) SendMsg(msg []byte) {
	if msg == nil {
		return
	}
	if c.config.ConnRWSeparate {
		c.Send <- msg
	} else {
		c.socketWriteMsg(msg)
	}
}

func (c *wsConnect) GetData() *types.ConnectData {
	return &types.ConnectData{
		UserId:        c.UserId,
		DeviceId:      c.DeviceId,
		Addr:          c.Addr,
		ServerIp: 	   c.ServerIp,
		RegisterTime:  c.RegisterTime,
		HeartbeatTime: c.HeartbeatTime,
		SetCacheTime:  c.SetCacheTime,
	}
}

func (c *wsConnect) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+uint64(globalkey.HeartbeatExpirationTime) <= currentTime {
		timeout = true
	}
	return
}

func (c *wsConnect) HeartBeat(req *types.ConnectSendMsg) (msg *types.ConnectSendResp) {
	//参数处理
	// 校验登入授权, 判断是否登入
	// 获取该连接缓存，设置心跳配置【HeartbeatTime:currentTime、IsLogoff：false】,可以通过go-zero内置lru实现【本地缓存-公共缓存redis】双缓存
	// 更新本地缓存的该连接心跳配置
	// 更新redis公共缓存的该连接心跳配置, 定时续期公共缓存
	CurrentTime := tool.GetUintNowTime()
	c.setHeartbeat(CurrentTime)
	if CurrentTime-c.SetCacheTime >= uint64(c.config.ConnectInfoTimeout/5*3) {
		cm.EventRenewal(&types.ConnectData{
			UserId:        c.UserId,
			DeviceId:      c.DeviceId,
		})
		c.setCacheTime(CurrentTime)
	}
	return &types.ConnectSendResp{
		Seq:  req.Seq,
		Cmd:  "heartbeat",
		Code: 200,
		Data: "ok",
	}
}

func (c *wsConnect) Ping(req *types.ConnectSendMsg) (msg *types.ConnectSendResp) {
	//参数处理
	// 校验登入授权, 判断是否登入
	// 算为一次心跳；so获取该连接缓存，设置心跳配置【HeartbeatTime:currentTime、IsLogoff：false】,可以通过go-zero内置lru实现【本地缓存-公共缓存redis】双缓存
	// 更新本地缓存的该连接心跳配置
	// 更新redis公共缓存的该连接心跳配置
	return &types.ConnectSendResp{
		Seq:  req.Seq,
		Cmd:  "ping",
		Code: 200,
		Data: "pong",
	}
}


func (c *wsConnect) socketWriteMsg(msg []byte) {
	// 每个连接一个协程的情况，这里的写消息可能会有并发竞争的情况，需要加锁
	c.SocketLock.Lock()
	err := c.Socket.WriteMessage(websocket.TextMessage, msg)
	c.SocketLock.Unlock()
	if err != nil {
		c.Errorf("%s 发送消息错误: %+v", c.Addr, err)
	}
}
func (c *wsConnect) setHeartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
}
func (c *wsConnect) setCacheTime(currentTime uint64) {
	c.SetCacheTime = currentTime
}
