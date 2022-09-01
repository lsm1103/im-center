package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"im-center/common/xerr"
	"im-center/service/connect/rpc/connect"
	"strconv"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOnlineUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineUsersLogic {
	return &GetOnlineUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  --------------------------------------------------------------------------------------------------------
func (l *GetOnlineUsersLogic) GetOnlineUsers(in *chat.GetUserConnectListReq) (resp *chat.UserConnectListResp, err error) {
	node := l.svcCtx.RpcU.GetNode()
	if node == nil {
		return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
	}
	list, err := node.GetUserConnectList(l.ctx, &connect.GetUserConnectListReq{
		UserId: strconv.FormatInt(in.UserId, 10),
		Offset: uint64(in.Offset),
		Limit: uint64(in.Limit),
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
