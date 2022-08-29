package service

import (
	"context"
	"github.com/jinzhu/copier"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

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
	err = copier.Copy(resp, l.svcCtx.RpcU.NodeList)
	if err != nil {
		return nil, err
	}
	return
}
