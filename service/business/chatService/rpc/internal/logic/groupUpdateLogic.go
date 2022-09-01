package logic

import (
	"context"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUpdateLogic {
	return &GroupUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupUpdateLogic) GroupUpdate(in *chat.GroupUpdateReq) (*chat.NullResp, error) {
	err := l.svcCtx.GroupModel.Update(nil, &database.Group{
		Id:         in.GroupType,
		Name:       in.Name,
		CreateUser: in.CreateUser,
		Ico:        in.Ico,
		Remark:     in.Remark,
		ParentId:   in.ParentId,
		GroupType:  in.GroupType,
		Rank:       in.Rank,
		Status:     in.Status,
	})
	if err != nil {
		return nil, err
	}

	return &chat.NullResp{}, nil
}
