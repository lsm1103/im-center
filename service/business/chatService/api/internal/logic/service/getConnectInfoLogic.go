package service

import (
	"context"
	"im-center/common/xerr"
	"im-center/service/connect/rpc/connect"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConnectInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConnectInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetConnectInfoLogic {
	return GetConnectInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConnectInfoLogic) GetConnectInfo(req types.ConnectUid) (resp *types.ConnectInfoResp, err error) {
	node := l.svcCtx.RpcU.GetNode()
	if node == nil {
		return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
	}
	info, err := node.GetConnectInfo(l.ctx, &connect.GetConnectInfoReq{
		UserId:   req.UserId,
		DeviceId: req.DeviceId,
	})
	if err != nil {
		return nil, err
	}
	return &types.ConnectInfoResp{
		UserId:         info.UserId,
		DeviceId:       info.DeviceId,
		ServerIp:       info.ServerIp,
		ConnectIp:      info.ConnectIp,
		RegisterTime:   info.RegisterTime,
		HeartbeatTime:  info.HeartbeatTime,
		UnRegisterTime: info.UnRegisterTime,
		DeviceInfo:     info.DeviceInfo,
	}, nil
}
