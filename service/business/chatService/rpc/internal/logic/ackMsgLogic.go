package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"im-center/common/globalkey"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"
)

type AckMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAckMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AckMsgLogic {
	return &AckMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AckMsgLogic) AckMsg(in *chat.AckMsgReq) (*chat.NullResp, error) {
	err := l.svcCtx.OfflineMsgModel.Update(nil, &database.OfflineMsg{
		UserId:     in.UserId,
		DeviceId:   in.DeviceId,
		ObjectType: globalkey.MsgType[in.ObjectType],
		ObjectId:   in.ObjectId,
		LastAckSeq: in.Seq,
	})
	if err != nil {
		return nil, err
	}
	return &chat.NullResp{}, nil
}
