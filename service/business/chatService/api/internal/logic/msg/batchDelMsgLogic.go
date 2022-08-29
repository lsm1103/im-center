package msg

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDelMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDelMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchDelMsgLogic {
	return BatchDelMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDelMsgLogic) BatchDelMsg(req types.BatchDelMsgReq) error {
	// todo: add your logic here and delete this line

	return nil
}
