package friend

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserExitGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserExitGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserExitGroupLogic {
	return UserExitGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserExitGroupLogic) UserExitGroup(req types.UserExitGroupReq) error {
	// todo: add your logic here and delete this line

	return nil
}
