package logic

import (
	"context"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFriendLogic) AddFriend(in *chat.AddFriendReq) (*chat.NullResp, error) {
	_, err := l.svcCtx.FriendModel.Insert(nil, &database.Friend{
		ApplyUser:   in.ApplyUser,
		ApplyDevice: in.ApplyDevice,
		AcceptUser:  in.AcceptUser,
		ApplyReason: in.ApplyReason,
	})
	if err != nil {
		return nil, err
	}
	return &chat.NullResp{}, nil
}
