package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im-center/common/globalkey"
	"im-center/common/tool"
	"im-center/common/uniqueid"
	"im-center/common/xerr"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/connect/rpc/connectclient"
	"im-center/service/model/database"
	"strconv"
)

type SendOneMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendOneMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendOneMsgLogic {
	return &SendOneMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendOneMsgLogic) SendOneMsg(in *chat.SendOneMsgReq) (*chat.NullResp, error) {
	if err := l.svcCtx.ModelHandle.Trans(func(session sqlx.Session) error {
		lastId := uniqueid.GenId()
		_, err := l.svcCtx.SingleMsgModel.Insert(session, &database.SingleMsg{
			Id: uniqueid.GenId(),
			Seq: lastId,
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
			ObjectType: globalkey.MsgType[globalkey.SingleMsg],
			ObjectId:   in.ReceiverId,
			NewestSeq:  lastId,
		})
		return err
	}); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "消息写入数据库失败，err:s%", err.Error())
	}

	_, err := l.svcCtx.RpcU.GetNode().SendOneMsg(l.ctx, &connectclient.SendOneMsgReq{
		SenderType:       strconv.FormatInt(in.SenderType, 10),
		SenderId:         strconv.FormatInt(in.SenderId, 10),
		SenderDeviceId:   in.SenderDeviceId,
		ReceiverId:       strconv.FormatInt(in.ReceiverUserId, 10),
		ReceiverDeviceId: in.ReceiverDeviceId,
		ParentId:         strconv.FormatInt(in.ParentId, 10),
		SendTime:         in.SendTime,
		MsgType: 		  strconv.FormatInt(in.MsgType, 10),
		MsgContent:       in.Content,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.USER_OPERATION_ERR), "消息写入数据库成功，发送消息失败，err:s%", err.Error())
	}
	return &chat.NullResp{}, nil
}
