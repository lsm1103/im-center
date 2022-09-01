package service

import (
	"context"
	"github.com/jinzhu/copier"
	"im-center/common/xerr"
	"im-center/service/connect/rpc/connect"
	"strconv"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

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
	node := l.svcCtx.RpcU.GetNode()
	if node == nil {
		return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
	}
	list, err := node.GetUserConnectList(l.ctx, &connect.GetUserConnectListReq{
		UserId: strconv.FormatInt(req.UserId, 10),
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return nil, err
	}

	err = copier.Copy(resp, list)
	if err != nil {
		return nil, err
	}

	return
}
