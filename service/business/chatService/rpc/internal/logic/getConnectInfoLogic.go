package logic

import (
	"context"
	"im-center/common/xerr"
	"im-center/service/connect/rpc/connect"
	"strconv"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConnectInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConnectInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConnectInfoLogic {
	return &GetConnectInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConnectInfoLogic) GetConnectInfo(in *chat.ConnectUid) (*chat.ConnectItem, error) {
	node := l.svcCtx.RpcU.GetNode()
	if node == nil {
		return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
	}
	info, err := node.GetConnectInfo(l.ctx, &connect.GetConnectInfoReq{
		UserId:   strconv.FormatInt(in.UserId, 10),
		DeviceId: in.DeviceId,
	})
	if err != nil {
		return nil, err
	}
	return &chat.ConnectItem{
		UserId:         in.UserId,
		DeviceId:       info.DeviceId,
		ServerIp:       info.ServerIp,
		ConnectIp:      info.ConnectIp,
		RegisterTime:   info.RegisterTime,
		HeartbeatTime:  info.HeartbeatTime,
		UnRegisterTime: info.UnRegisterTime,
		DeviceInfo:     info.DeviceInfo,
	}, nil
}
