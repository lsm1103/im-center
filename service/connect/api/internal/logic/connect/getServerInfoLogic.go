package connect

import (
	"context"
	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetServerInfoLogic {
	return GetServerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerInfoLogic) GetServerInfo(req types.GetServerInfoReq) (resp *types.ServerInfoResp, err error) {
	info, err := l.svcCtx.Cs.GetServiceInfo(nil)
	if err != nil {
		return nil, err
	}
	return info, nil
}
