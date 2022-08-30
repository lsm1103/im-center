package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncMsgLogic {
	return &SyncMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncMsgLogic) SyncMsg(in *chat.SyncMsgReq) (*chat.SyncMsgResp, error) {
	// todo: add your logic here and delete this line

	return &chat.SyncMsgResp{}, nil
}
