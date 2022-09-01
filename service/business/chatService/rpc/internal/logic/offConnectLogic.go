package logic

import (
	"context"
	"im-center/common/xerr"
	"im-center/service/connect/rpc/connect"
	"strconv"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type OffConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOffConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OffConnectLogic {
	return &OffConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OffConnectLogic) OffConnect(in *chat.ConnectUid) (resp *chat.NullResp, err error) {
	node := l.svcCtx.RpcU.GetNode()
	if node == nil {
		return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
	}
	_, err = node.OffConnect(l.ctx, &connect.OffConnectReq{
		UserId: strconv.FormatInt(in.UserId, 10),
		DeviceId: in.DeviceId,
	})
	if err != nil {
		return nil, err
	}
	return
}
