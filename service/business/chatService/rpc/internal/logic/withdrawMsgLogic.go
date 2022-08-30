package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawMsgLogic {
	return &WithdrawMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawMsgLogic) WithdrawMsg(in *chat.WithdrawMsgReq) (*chat.NullResp, error) {
	// todo: add your logic here and delete this line

	return &chat.NullResp{}, nil
}
