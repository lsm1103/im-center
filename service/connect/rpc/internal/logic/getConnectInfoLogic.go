package logic

import (
	"context"
	"im-center/service/connect/internal/types"
	"strconv"

	"im-center/service/connect/rpc/connect"
	"im-center/service/connect/rpc/internal/svc"

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

func (l *GetConnectInfoLogic) GetConnectInfo(in *connect.GetConnectInfoReq) (*connect.ConnectInfo, error) {
	info, err := l.svcCtx.CS.GetConnectInfo(&types.ConnectUid{
		UserId:   in.UserId,
		DeviceId: in.DeviceId,
	})
	if err != nil {
		return nil, err
	}

	return &connect.ConnectInfo{
		UserId:         info.UserId,
		DeviceId:       info.DeviceInfo,
		ServerIp:       info.ServerIp,
		ConnectIp:      info.ConnectIp,
		RegisterTime:   strconv.FormatUint(info.RegisterTime, 10),
		HeartbeatTime:  strconv.FormatUint(info.HeartbeatTime, 10),
		UnRegisterTime: strconv.FormatUint(info.UnRegisterTime, 10),
		DeviceInfo:     info.DeviceInfo,
	}, nil
}
