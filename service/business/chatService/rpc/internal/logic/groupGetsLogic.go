package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"im-center/common/tool"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupGetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupGetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupGetsLogic {
	return &GroupGetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupGetsLogic) GroupGets(in *chat.GetsReq) (resp *chat.GroupGetsResp, err error) {
	findReq := &tool.GetsReq{
		OrderBy:  in.OrderBy,
		Sort:     in.Sort,
		Current:  in.Current,
		PageSize: in.PageSize,
	}
	for _,item := range in.Query{
		findReq.Query = append(findReq.Query, &tool.GetsQueryItem{
			Key:        item.Key,
			Val:        item.Val,
			Handle:     item.Handle,
			NextHandle: item.NextHandle,
		})
	}
	all, err := l.svcCtx.GroupModel.FindAll(findReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, all)
	if err != nil {
		return nil, err
	}
	return
}
