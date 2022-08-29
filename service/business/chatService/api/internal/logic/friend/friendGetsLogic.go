package friend

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendGetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendGetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) FriendGetsLogic {
	return FriendGetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendGetsLogic) FriendGets(req types.GetsReq) (resp *types.FriendGetsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
