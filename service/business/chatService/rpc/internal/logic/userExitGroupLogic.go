package logic

import (
	"context"
	"im-center/common/globalkey"
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
	err := l.svcCtx.UserGroupModel.SoftDelete(nil, &database.UserGroup{
		UserId:  in.UserId,
		GroupId: in.GroupId,
		Status:  globalkey.Del,
	})
	if err != nil {
		return nil, err
	}
	return &chat.NullResp{}, nil
}
