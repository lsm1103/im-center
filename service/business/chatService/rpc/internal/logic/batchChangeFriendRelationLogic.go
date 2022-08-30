package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchChangeFriendRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchChangeFriendRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchChangeFriendRelationLogic {
	return &BatchChangeFriendRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchChangeFriendRelationLogic) BatchChangeFriendRelation(in *chat.BatchChangeFriendRelationReq) (*chat.NullResp, error) {
	// todo: add your logic here and delete this line

	return &chat.NullResp{}, nil
}
