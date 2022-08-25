package connect

import (
	"context"

	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"

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
	l.Infof("api SendManyMsg, %+v", req)
	err = l.svcCtx.Cs.SendManyMsg(&req)
	if err != nil {
		return nil, err
	}
	return
}
