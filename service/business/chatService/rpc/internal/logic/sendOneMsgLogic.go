package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &chat.NullResp{}, nil
}
