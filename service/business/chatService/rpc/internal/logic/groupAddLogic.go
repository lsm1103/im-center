package logic

import (
	"context"
	"im-center/common/globalkey"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupAddLogic {
	return &GroupAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupAddLogic) GroupAdd(in *chat.GroupAddReq) (*chat.NullResp, error) {
	// 通过群聊类型，用不同的策略来验证本次创建的群聊是否合法
	switch in.GroupType {
	case globalkey.Department:
	case globalkey.UserGroup:
	case globalkey.Group:
	case globalkey.Circle:
	case globalkey.Topics:
	}
	// 写入数据库，日志处理由网关来发起异步任务
	_, err := l.svcCtx.GroupModel.Insert(nil, &database.Group{
		Name:       in.Name,
		CreateUser: in.CreateUser,
		Ico:        in.Ico,
		Remark:     in.Remark,
		ParentId:   in.ParentId,
		GroupType:  in.GroupType,
		Rank:       in.Rank,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
