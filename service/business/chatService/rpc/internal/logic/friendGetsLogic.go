package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendGetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendGetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendGetsLogic {
	return &FriendGetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendGetsLogic) FriendGets(in *chat.GetsReq) (*chat.FriendGetsResp, error) {
	// todo: add your logic here and delete this line

	return &chat.FriendGetsResp{}, nil
}
