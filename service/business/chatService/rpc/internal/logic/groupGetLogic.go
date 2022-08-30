package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupGetLogic {
	return &GroupGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  --------------------------------------------------------------------------------------------------------
func (l *GroupGetLogic) GroupGet(in *chat.GroupGetReq) (*chat.GroupItem, error) {
	// todo: add your logic here and delete this line

	return &chat.GroupItem{}, nil
}
