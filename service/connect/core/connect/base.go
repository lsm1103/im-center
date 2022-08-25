package connect

import (
	"im-center/service/connect/internal/types"
)

type Connect interface {
	// 每一个连接开启一个协程，处理读消息操作，写消息通过http/rpc请求的协程来发送；就是说如果要减少协程数量，就把每个连接的写消息队列操作取消
	Read()
	// 向连接写数据
	Write()
	// 判断心跳是否超时
	IsHeartbeatTimeout(currentTime uint64) (timeout bool)
	// 主动关闭
	ActiveClose()
	// 被动关闭
	PassiveClose()
	// 发消息
	SendMsg(msg []byte)
	// 获取数据
	GetData() *types.ConnectData
	// 心跳
	HeartBeat(req *types.ConnectSendMsg) (msg *types.ConnectSendResp)
	// ping
	Ping(req *types.ConnectSendMsg) (msg *types.ConnectSendResp)
}