package msg

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) SyncMsgLogic {
	return SyncMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncMsgLogic) SyncMsg(req types.SyncMsgReq) (resp *types.SyncMsgResp, err error) {
	// todo: add your logic here and delete this line

	return
}
