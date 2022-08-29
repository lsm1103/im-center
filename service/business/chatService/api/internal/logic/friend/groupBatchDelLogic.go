package friend

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupBatchDelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupBatchDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupBatchDelLogic {
	return GroupBatchDelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupBatchDelLogic) GroupBatchDel(req types.GroupBatchDelReq) error {
	// todo: add your logic here and delete this line

	return nil
}
