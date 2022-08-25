package tcp

import (
	"context"
	"im-center/service/connect/internal/config"

	"github.com/alberliu/gn"
	"github.com/zeromicro/go-zero/core/logx"

	"im-center/common/globalkey"
	"im-center/common/tool"
	"im-center/service/connect/core/connect"
	"im-center/service/connect/core/connectManager"
	"im-center/service/connect/internal/types"
)

type tcpConnect struct {
	logx.Logger
	ctx    		  context.Context
	config        *config.Config

	Socket        *gn.Conn 		  // 连接
	Send          chan []byte     // 待发送的数据
	Addr          string          // 客户端地址
	UserId        string          // 用户Id
	DeviceId      string          // 登录的平台Id app/web/ios
	RegisterTime  uint64          // 用户上次上线时间
	HeartbeatTime uint64          // 用户上次心跳时间
	SetCacheTime  uint64          // 用户上次设置缓存时间
	ServerIp	  string          // 服务器ip
	//UnRegisterTime uint64       // 用户上次下线时间
	//DeviceInfo    string		  // 用户设备信息
}

func NewTcpConnect(ctx context.Context, config *config.Config, in types.TcpConnect) connect.Connect {
	return &tcpConnect{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		config: config,

		Socket:        in.Socket,
		Send:          make(chan []byte, 100),
		UserId:        in.UserId,
		DeviceId:      in.DeviceId,
		Addr:          in.Addr,
		RegisterTime:  in.CurrentTime,
		HeartbeatTime: in.CurrentTime,
		SetCacheTime:  in.CurrentTime,
		ServerIp: 	   config.GetServerIp(),
	}
}

func (c *tcpConnect) Read() {
	c.Info("implement me")
}

func (c *tcpConnect) Write() {
	defer func() {
		c.Errorf("defer, 关闭Write协程 %+v", c)
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

func (c *tcpConnect) ActiveClose() {
	c.Socket.Close()
}
func (c *tcpConnect) PassiveClose() {
	close(c.Send)
}

func (c *tcpConnect) SendMsg(msg []byte) {
	if msg == nil {
		return
	}
	if c.config.ConnRWSeparate {
		c.Send <- msg
	} else {
		c.socketWriteMsg(msg)
	}
}

func (c *tcpConnect) GetData() *types.ConnectData {
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

func (c *tcpConnect) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+uint64(globalkey.HeartbeatExpirationTime) <= currentTime {
		timeout = true
	}
	return
}

func (c *tcpConnect) HeartBeat(req *types.ConnectSendMsg) (msg *types.ConnectSendResp) {
	//参数处理
	// 校验登入授权, 判断是否登入
	// 获取该连接缓存，设置心跳配置【HeartbeatTime:currentTime、IsLogoff：false】,可以通过go-zero内置lru实现【本地缓存-公共缓存redis】双缓存
	// 更新本地缓存的该连接心跳配置
	// 更新redis公共缓存的该连接心跳配置, 定时续期公共缓存
	CurrentTime := tool.GetUintNowTime()
	c.setHeartbeat(CurrentTime)
	if CurrentTime-c.SetCacheTime >= uint64(c.config.ConnectInfoTimeout/5*3) {
		connectManager.GetCM().EventRenewal(&types.ConnectData{
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

func (c *tcpConnect) Ping(req *types.ConnectSendMsg) (msg *types.ConnectSendResp) {
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


func (c *tcpConnect) socketWriteMsg(msg []byte) {
	_, err := c.Socket.Write(msg)
	if err != nil {
		c.Errorf("%s 发送消息错误: %+v", c.Addr, err)
	}
}
func (c *tcpConnect) setHeartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
}
func (c *tcpConnect) setCacheTime(currentTime uint64) {
	c.SetCacheTime = currentTime
}

