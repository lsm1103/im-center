package logic

import (
	"context"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserExitGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserExitGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserExitGroupLogic {
	return &UserExitGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserExitGroupLogic) UserExitGroup(in *chat.UserExitGroupReq) (*chat.NullResp, error) {
	// todo 修改updata，实现多字段确定唯一问题
	err := l.svcCtx.UserGroupModel.SoftDelete(nil, &database.UserGroup{
		UserId:  in.UserId,
		GroupId: in.GroupId,
		Status:  -1,
	})
	if err != nil {
		return nil, err
	}
	return &chat.NullResp{}, nil
}
