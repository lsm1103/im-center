package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAddGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAddGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAddGroupLogic {
	return &UserAddGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  用户群组关系
func (l *UserAddGroupLogic) UserAddGroup(in *chat.UserAddGroupReq) (*chat.NullResp, error) {
	// todo: add your logic here and delete this line

	return &chat.NullResp{}, nil
}
