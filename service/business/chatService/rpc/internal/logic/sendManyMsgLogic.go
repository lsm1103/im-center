package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im-center/common/globalkey"
	"im-center/common/tool"
	"im-center/common/uniqueid"
	"im-center/common/xerr"
	"im-center/service/connect/rpc/connectclient"
	"im-center/service/model/database"
	"strconv"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendManyMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendManyMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendManyMsgLogic {
	return &SendManyMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendManyMsgLogic) SendManyMsg(in *chat.SendManyMsgReq) (*chat.NullResp, error) {
	// 异步定时任务：如果当前seq和最新的seq对不上，说明该 用户-设备-消息房间 有离线消息

	// 查询该群组
	all, err := l.svcCtx.UserGroupModel.FindAll(&tool.GetsReq{
		Query: []*tool.GetsQueryItem{
			{
				Key:    "group_id",
				Val:    strconv.FormatInt(in.ReceiverId, 10),
				Handle: "=",
			},
		},
		Current:  0,
		PageSize: 100000,
	})
	if err != nil {
		return nil, err
	}
	// 存储消息到数据库
	if err := l.svcCtx.ModelHandle.Trans(func(session sqlx.Session) error {
		lastId := uniqueid.GenId()
		_, err := l.svcCtx.GroupMsgModel.Insert(session, &database.GroupMsg{
			Id: uniqueid.GenId(),
			SenderType:       in.SenderType,
			SenderId:         in.SenderId,
			SenderDeviceId:   in.SenderDeviceId,
			ReceiverId:       in.ReceiverId,
			ReceiverDeviceId: in.ReceiverDeviceId,
			MsgType:          in.MsgType,
			Content:          in.Content,
			ParentId:         in.ParentId,
			SendTime:         tool.StrToTime(in.SendTime),
		})
		if err != nil {
			return err
		}
		_, err = l.svcCtx.OfflineMsgModel.Insert(session, &database.OfflineMsg{
			Id: uniqueid.GenId(),
			UserId:     in.SenderId,
			DeviceId:   in.SenderDeviceId,
			ObjectType: globalkey.MsgType[globalkey.GroupMsg],
			ObjectId:   in.ReceiverId,
			NewestSeq:  lastId,
		})
		return err
	}); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "消息写入数据库失败，err:s%", err.Error())
	}
	// 发送消息，todo 可以这这一步优化成查询各在线节点的在线用户，进行按节点分组后，并发发给相对应的节点，避免节点间的转发
	receiverIds := []string{}
	for _,item := range all{
		receiverIds = append(receiverIds, strconv.FormatInt(item.Id, 10))
	}
	_, err = l.svcCtx.RpcU.GetNode().SendManyMsg(l.ctx, &connectclient.SendManyMsgReq{
		SenderType:           strconv.FormatInt(in.SenderType, 10),
		SenderId:             strconv.FormatInt(in.SenderId, 10),
		SenderDeviceId:       in.SenderDeviceId,
		ReceiverId:           receiverIds,
		ReceiverDeviceId:     in.ReceiverDeviceId,
		ParentId:             strconv.FormatInt(in.ParentId, 10),
		SendTime:             in.SendTime,
		At:                   nil,
		MsgType:              strconv.FormatInt(in.MsgType, 10),
		MsgContent:           in.Content,
	})
	if err != nil {
		return nil, err
	}
	return &chat.NullResp{}, nil
}
