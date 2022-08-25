package connect

import (
	"im-center/service/connect/internal/types"
	"im-center/service/connect/internal/utils"
)

func ProcessData(c Connect, message []byte) {
	//参数处理
	req := utils.ParseConnectMsg(message)
	if req == nil {
		c.SendMsg([]byte(`{"seq":"-", "cmd":"error", "code":500, "data":"数据不合法"}`) )
		return
	}
	// 按cmd处理
	msgBody := &types.ConnectSendResp{}
	switch req.Cmd {
	case "heartbeat":
		msgBody = c.HeartBeat(req)
	case "ping":
		msgBody = c.Ping(req)
	default:
		//todo 修改该用户连接为异常连接
		msgBody.Seq = req.Seq
		msgBody.Cmd = "error"
		msgBody.Code = 405
		msgBody.Data = "cmd类型不支持"
	}
	//发送消息
	msg := utils.BuildConnectMsg(msgBody)

	if msg != nil {
		c.SendMsg(msg)
	}
}