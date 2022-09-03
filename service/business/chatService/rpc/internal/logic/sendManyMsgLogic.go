package logic

import (
	"context"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendManyMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendManyMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendManyMsgLogic {
	return &SendManyMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendManyMsgLogic) SendManyMsg(in *chat.SendManyMsgReq) (*chat.NullResp, error) {
	// todo: 1、写入消息表；
	// 2、更新各用户-设备-消息房间的最新seq；
	// 3、尝试下发给各用户-设备-消息房间；
	// 4、下发成功更新各用户-设备-消息房间的当前seq
	// 5、如果当前seq和最新的seq对不上，说明该 用户-设备-消息房间 有离线消息

	return &chat.NullResp{}, nil
}
