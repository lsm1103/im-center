package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDelMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDelMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDelMsgLogic {
	return &BatchDelMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchDelMsgLogic) BatchDelMsg(in *chat.BatchDelMsgReq) (*chat.NullResp, error) {
	// todo: add your logic here and delete this line

	return &chat.NullResp{}, nil
}
