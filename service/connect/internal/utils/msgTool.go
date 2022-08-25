package utils

import (
	"encoding/json"
	"fmt"
	"im-center/service/connect/internal/types"
	"im-center/service/connect/rpc/connect"
	"io"
)

// 构建消息包
func BuildConnectMsg(msgBody *types.ConnectSendResp) []byte {
	msg, err := json.Marshal(msgBody)
	if err != nil {
		fmt.Println("构建消息包失败:", err)
		return nil
	}
	return msg
}

// 解析消息包
func ParseConnectMsg(msg []byte) (req *types.ConnectSendMsg)  {
	req = &types.ConnectSendMsg{}
	err := json.Unmarshal(msg, req)
	if err != nil {
		fmt.Println("解析消息包失败:", err)
		return nil
	}
	return
}

func BuildMsg (msg *types.SendOneMsgReq) []byte {
	tmp := types.Msg{
		Seq:  "124214124",
		Cmd:  "msg",
		Data: types.Data{
			SenderType:     msg.SenderType,
			SenderId:       msg.SenderId,
			SenderDeviceId: msg.SenderDeviceId,
			ParentId:       msg.ParentId,
			SendTime:       msg.SendTime,
			MsgType:        msg.MsgType,
			//At:             msg.At,
			MsgContent: 	msg.MsgContent,
		},
	}

	//switch msg.MsgType {
	//case "text":
	//	tmp.Data.TextMsg = msg.TextMsg
	//case "img":
	//	tmp.Data.ImgMsg = msg.ImgMsg
	//case "audio":
	//	tmp.Data.AudioMsg = msg.AudioMsg
	//case "video":
	//	tmp.Data.VideoMsg = msg.VideoMsg
	//case "link":
	//	tmp.Data.LinkMsg = msg.LinkMsg
	//case "markdown":
	//	tmp.Data.MarkdownMsg = msg.MarkdownMsg
	//case "actionCard":
	//	tmp.Data.ActionCardMsg = msg.ActionCardMsg
	//case "feedCard":
	//	tmp.Data.FeedCardMsg = msg.FeedCardMsg
	//case "file":
	//	tmp.Data.MsgContent = msg.MsgContent.(types.File)
	//case "position":	//位置
	//	tmp.Data.MsgContent = msg.MsgContent.(types.Position)
	//case "notice":	// 通知
	//	tmp.Data.MsgContent = msg.MsgContent.(types.Notice)
	//default:
	//	fmt.Println("该消息类型不支持:", msg.MsgType)
	//	return nil
	//}
	r, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println("构建消息包失败:", err)
		return nil
	}
	return []byte(r)
}
func ParseMsg(msg io.ReadCloser) (req *types.Msg, err error)  {
	req = &types.Msg{}
	err = json.NewDecoder(msg).Decode(req)
	if err != nil {
		fmt.Println("解析消息包失败:", err)
	}
	return
}


// 解析消息包
func ParseOneMsg(msg io.ReadCloser) (req *types.SendOneMsgReq, err error)  {
	req = &types.SendOneMsgReq{}
	err = json.NewDecoder(msg).Decode(req)
	if err != nil {
		fmt.Println("解析消息包失败:", err)
	}
	return
}
func ParseMassMsg(msg io.ReadCloser) (req *types.SendManyMsgReq, err error)  {
	req = &types.SendManyMsgReq{}
	err = json.NewDecoder(msg).Decode(req)
	if err != nil {
		fmt.Println("解析消息包失败:", err)
	}
	return
}


func BuildRpcOneMsg(in *connect.SendOneMsgReq) (req *types.SendOneMsgReq, err error){
	req = &types.SendOneMsgReq{
		SenderType:       in.SenderType,
		SenderId:         in.SenderId,
		SenderDeviceId:   in.SenderDeviceId,
		ReceiverId:       in.ReceiverId,
		ReceiverDeviceId: in.ReceiverDeviceId,
		ParentId:         in.ParentId,
		SendTime:         in.SendTime,
		MsgType:          in.MsgType,
	}
	switch in.MsgType {
	case "text":
		var msgCenter types.Text
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "img":
		var msgCenter types.Img
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "audio":
		var msgCenter types.Audio
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "video":
		var msgCenter types.Video
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "link":
		var msgCenter types.Link
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "markdown":
		var msgCenter types.Markdown
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "actionCard":
		var msgCenter types.ActionCard
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "feedCard":
		var msgCenter types.FeedCard
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	}
	return
}

func BuildRpcManyMsg(in *connect.SendManyMsgReq) (req *types.SendManyMsgReq, err error){
	req = &types.SendManyMsgReq{
		SenderType:       in.SenderType,
		SenderId:         in.SenderId,
		SenderDeviceId:   in.SenderDeviceId,
		ReceiverId:       in.ReceiverId,
		ReceiverDeviceId: in.ReceiverDeviceId,
		ParentId:         in.ParentId,
		SendTime:         in.SendTime,
		At:               types.At{
			IsAtAll:   in.At.IsAtAll,
			AtUserIds: in.At.AtUserIds,
		},
		MsgType:          in.MsgType,
	}
	switch in.MsgType {
	case "text":
		var msgCenter types.Text
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "img":
		var msgCenter types.Img
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "audio":
		var msgCenter types.Audio
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "video":
		var msgCenter types.Video
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "link":
		var msgCenter types.Link
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "markdown":
		var msgCenter types.Markdown
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "actionCard":
		var msgCenter types.ActionCard
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	case "feedCard":
		var msgCenter types.FeedCard
		err := json.Unmarshal([]byte(in.MsgContent), &msgCenter)
		if err != nil {
			return nil, err
		}
		req.MsgContent = msgCenter
	}
	return
}