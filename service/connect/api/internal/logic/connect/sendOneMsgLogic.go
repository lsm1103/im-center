package connect

import (
	"context"
	"im-center/common/xerr"
	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"

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
	l.Infof("api SendOneMsg, %+v", req)
	err = l.svcCtx.Cs.SendOneMsg(&req, false)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.USER_OPERATION_ERR, err.Error())
	}
	return
}
