package service

import (
	"context"
	"im-center/common/xerr"
	"im-center/service/connect/rpc/connect"
	"strconv"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OffConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOffConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) OffConnectLogic {
	return OffConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OffConnectLogic) OffConnect(req types.ConnectUid) (resp *types.NullResp, err error) {
	node := l.svcCtx.RpcU.GetNode()
	if node == nil {
		return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
	}
	_, err = node.OffConnect(l.ctx, &connect.OffConnectReq{
		UserId: strconv.FormatInt(req.UserId, 10),
		DeviceId: req.DeviceId,
	})
	if err != nil {
		return nil, err
	}
	return
}
