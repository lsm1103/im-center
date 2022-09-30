package logic

import (
	"context"
	"im-center/common/tool"
	"strconv"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncMsgLogic {
	return &SyncMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncMsgLogic) SyncMsg(in *chat.SyncMsgReq) (resp *chat.SyncMsgResp, err error) {
	for _,single := range in.SingleFilters {
		singleAll, err := l.svcCtx.SingleMsgModel.FindAll(&tool.GetsReq{
			Query: []*tool.GetsQueryItem{
				{
					Key:        "seq",
					Val: 		strconv.FormatInt(single.Seq, 10),
					Handle:     ">",
					NextHandle: "and",
				},
				{
					Key:        "receiver_id",
					Val:        strconv.FormatInt(single.FriendId, 10),
					Handle:     "=",
					NextHandle: "and",
				},
				{
					Key:        "receiver_device_id",
					Val:        in.DeviceId,
					Handle:     "=",
					NextHandle: "and",
				},
			},
			Current:  0,
			PageSize: 300,
		})
		if err != nil {
			return nil, err
		}
		msgs := chat.SingleMsgs{}
		for _,msg := range singleAll{
			msgs.Msgs = append(msgs.Msgs, &chat.SingleMsgItem{
				Id:                   msg.Id,
				Seq: 				  msg.Seq,
				SenderType:           strconv.FormatInt(msg.SenderType, 10),
				SenderId:             msg.SenderId,
				SenderDeviceId:       msg.SenderDeviceId,
				ReceiverId:           msg.ReceiverId,
				MsgType:              strconv.FormatInt(msg.MsgType, 10),
				Content:              msg.Content,
				ParentId:             msg.ParentId,
				SendTime:             tool.FmtTime(msg.SendTime),
				Status:               msg.Status,
				CreateTime:           tool.FmtTime(msg.CreateTime),
				UpdateTime:           tool.FmtTime(msg.UpdateTime),
			})
		}
		resp.SingleMsgList[strconv.FormatInt(single.FriendId, 10)] = &msgs
	}

	for _,group := range in.GroupFilters {
		groupAll, err := l.svcCtx.GroupMsgModel.FindAll(&tool.GetsReq{
			Query: []*tool.GetsQueryItem{
				{
					Key:        "seq",
					Val: 		strconv.FormatInt(group.Seq, 10),
					Handle:     ">",
					NextHandle: "and",
				},
				{
					Key:        "receiver_id",
					Val: 		strconv.FormatInt(group.GroupId, 10),
					Handle:     "=",
					NextHandle: "and",
				},
				{
					Key:        "receiver_device_id",
					Val:        in.DeviceId,
					Handle:     "=",
					NextHandle: "and",
				},
			},
			Current:  0,
			PageSize: 300,
		})
		if err != nil {
			return nil, err
		}
		msgs := chat.GroupMsgs{}
		for _,msg := range groupAll{
			msgs.Msgs = append(msgs.Msgs, &chat.GroupMsgItem{
				Id:                   msg.Id,
				Seq: 				  msg.Seq,
				SenderType:           strconv.FormatInt(msg.SenderType, 10),
				SenderId:             msg.SenderId,
				SenderDeviceId:       msg.SenderDeviceId,
				ReceiverId:           msg.ReceiverId,
				MsgType:              strconv.FormatInt(msg.MsgType, 10),
				Content:              msg.Content,
				ParentId:             msg.ParentId,
				SendTime:             tool.FmtTime(msg.SendTime),
				Status:               msg.Status,
				CreateTime:           tool.FmtTime(msg.CreateTime),
				UpdateTime:           tool.FmtTime(msg.UpdateTime),
			})
		}
		resp.GroupMsgList[strconv.FormatInt(group.GroupId, 10)] = &msgs
	}

	return
}
