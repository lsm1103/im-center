package friend

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchChangeFriendRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchChangeFriendRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchChangeFriendRelationLogic {
	return BatchChangeFriendRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchChangeFriendRelationLogic) BatchChangeFriendRelation(req types.BatchChangeFriendRelationReq) error {
	// todo: add your logic here and delete this line

	return nil
}
