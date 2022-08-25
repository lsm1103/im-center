package connect

import (
	"context"

	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"

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
	resp, err = l.svcCtx.Cs.GetConnectInfo(&req)
	if err != nil {
		return nil, err
	}
	return
}
