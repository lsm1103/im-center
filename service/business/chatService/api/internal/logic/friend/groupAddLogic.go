package friend

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupAddLogic {
	return GroupAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupAddLogic) GroupAdd(req types.GroupAddReq) error {
	// todo: add your logic here and delete this line

	return nil
}
