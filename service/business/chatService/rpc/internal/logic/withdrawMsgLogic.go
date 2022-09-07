package logic

import (
	"context"
	"im-center/common/globalkey"
	"im-center/common/xerr"
	"im-center/service/model/database"

	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawMsgLogic {
	return &WithdrawMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawMsgLogic) WithdrawMsg(in *chat.WithdrawMsgReq) (*chat.NullResp, error) {
	switch in.ObjectType {
	case globalkey.SingleMsg:
		err := l.svcCtx.SingleMsgModel.SoftDelete(nil, &database.SingleMsg{
			Id:               in.Seq,
			Status:           globalkey.Withdraw,
		})
		if err != nil {
			return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
		}
	case globalkey.GroupMsg:
		err := l.svcCtx.GroupMsgModel.SoftDelete(nil, &database.GroupMsg{
			Id:               in.Seq,
			Status:           globalkey.Withdraw,
		})
		if err != nil {
			return nil, xerr.NewErrCode(xerr.USER_OPERATION_ERR)
		}
	}
	return &chat.NullResp{}, nil
}
