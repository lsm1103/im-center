package logic

import (
	"context"
	"im-center/service/connect/internal/types"

	"im-center/service/connect/rpc/connect"
	"im-center/service/connect/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type OffConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOffConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OffConnectLogic {
	return &OffConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OffConnectLogic) OffConnect(in *connect.OffConnectReq) (*connect.NullResp, error) {
	err := l.svcCtx.CS.OffConnect(&types.ConnectUid{
		UserId:   in.UserId,
		DeviceId: in.DeviceId,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
