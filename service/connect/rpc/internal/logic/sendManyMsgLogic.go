package logic

import (
	"context"
	"im-center/service/connect/internal/utils"

	"im-center/service/connect/rpc/connect"
	"im-center/service/connect/rpc/internal/svc"

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

func (l *SendManyMsgLogic) SendManyMsg(in *connect.SendManyMsgReq) (*connect.NullResp, error) {
	msg, err := utils.BuildRpcManyMsg(in)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.CS.SendManyMsg(msg)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
