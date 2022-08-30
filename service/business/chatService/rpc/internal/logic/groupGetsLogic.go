package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupGetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupGetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupGetsLogic {
	return &GroupGetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupGetsLogic) GroupGets(in *chat.GetsReq) (*chat.GroupGetsResp, error) {
	// todo: add your logic here and delete this line

	return &chat.GroupGetsResp{}, nil
}
