package logic

import (
	"context"
	"im-center/common/uniqueid"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"

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
	_, err := l.svcCtx.UserGroupModel.Insert(nil, &database.UserGroup{
		Id:         uniqueid.GenId(),
		UserId:     in.UserId,
		GroupId:    in.GroupId,
	})
	if err != nil {
		return nil, err
	}
	return &chat.NullResp{}, nil
}
