package logic

import (
	"context"
	"im-center/service/connect/internal/utils"

	"im-center/service/connect/rpc/connect"
	"im-center/service/connect/rpc/internal/svc"

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

func (l *SendOneMsgLogic) SendOneMsg(in *connect.SendOneMsgReq) (*connect.NullResp, error) {
	l.Infof("SendOneMsg: %s", in.String())
	msg, err := utils.BuildRpcOneMsg(in)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.CS.SendOneMsg(msg, in.IsLocal)
	if err != nil {
		return nil, err
	}
	return &connect.NullResp{}, nil
}
