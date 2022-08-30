package msg

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) WithdrawMsgLogic {
	return WithdrawMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WithdrawMsgLogic) WithdrawMsg(req types.WithdrawMsgReq) error {

	return nil
}
