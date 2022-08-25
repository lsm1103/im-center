package connect

import (
	"context"

	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConnectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserConnectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserConnectListLogic {
	return GetUserConnectListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserConnectListLogic) GetUserConnectList(req types.GetUserConnectListReq) (resp *types.UserConnectListResp, err error) {
	list, err := l.svcCtx.Cs.GetUserConnectList(&req)
	if err != nil {
		return nil, err
	}

	return &types.UserConnectListResp{
		UserConnectList: list.ConnectList,
	}, nil
}
