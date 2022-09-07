package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"im-center/common/tool"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendGetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendGetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendGetsLogic {
	return &FriendGetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendGetsLogic) FriendGets(in *chat.GetsReq) (resp *chat.GroupGetsResp, err error) {
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
	all, err := l.svcCtx.FriendModel.FindAll(findReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, all)
	if err != nil {
		return nil, err
	}
	return
}
