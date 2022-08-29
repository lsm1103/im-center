package friend

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAddGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAddGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserAddGroupLogic {
	return UserAddGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAddGroupLogic) UserAddGroup(req types.UserAddGroupReq) error {
	// todo: add your logic here and delete this line

	return nil
}
