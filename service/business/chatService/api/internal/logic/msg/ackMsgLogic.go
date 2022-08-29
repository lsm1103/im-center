package msg

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AckMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAckMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) AckMsgLogic {
	return AckMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AckMsgLogic) AckMsg(req types.AckMsgReq) (resp *types.NullResp, err error) {
	// todo: add your logic here and delete this line

	return
}
