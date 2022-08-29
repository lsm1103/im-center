package friend

import (
	"context"

	"im-center/service/business/chatService/api/internal/svc"
	"im-center/service/business/chatService/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupGetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupGetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupGetsLogic {
	return GroupGetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupGetsLogic) GroupGets(req types.GetsReq) (resp *types.GroupGetsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
