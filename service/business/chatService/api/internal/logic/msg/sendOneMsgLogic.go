package msg

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendOneMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendOneMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) SendOneMsgLogic {
	return SendOneMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendOneMsgLogic) SendOneMsg(req types.SendOneMsgReq) (resp *types.NullResp, err error) {
	// todo: add your logic here and delete this line

	return
}
