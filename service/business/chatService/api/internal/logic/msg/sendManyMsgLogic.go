package msg

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendManyMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendManyMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) SendManyMsgLogic {
	return SendManyMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendManyMsgLogic) SendManyMsg(req types.SendManyMsgReq) (resp *types.NullResp, err error) {
	// todo: add your logic here and delete this line

	return
}
