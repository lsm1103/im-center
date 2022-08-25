package connect

import (
	"context"

	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"

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
	l.Infof("api OffConnect, %+v", req)
	err = l.svcCtx.Cs.OffConnect(&req)
	if err != nil {
		return nil, err
	}
	return
}
