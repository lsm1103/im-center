package logic

import (
	"context"
	"im-center/common/tool"

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
	one, err := l.svcCtx.GroupModel.FindOne(in.GroupId)
	if err != nil {
		return nil, err
	}
	return &chat.GroupItem{
		Id:                   one.Id,
		Name:                 one.Name,
		CreateUser:           one.CreateUser,
		Ico:                  one.Ico,
		Remark:               one.Remark,
		ParentId:             one.ParentId,
		GroupType:            one.GroupType,
		Rank:                 one.Rank,
		Status:               one.Status,
		CreateTime:           tool.FmtTime(one.CreateTime),
		UpdateTime:           tool.FmtTime(one.UpdateTime),
	}, nil
}
