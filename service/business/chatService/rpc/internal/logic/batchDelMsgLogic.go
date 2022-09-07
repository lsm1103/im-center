package logic

import (
	"context"
	"fmt"
	"im-center/common/globalkey"
	"im-center/common/xerr"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"

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
	if _,ok := globalkey.MsgType[in.ObjectType]; !ok{
		return nil, xerr.NewErrMsg("找不到该消息对象类型")
	}
	success := []int64{}
	for _,seq := range in.Seqs {
		switch in.ObjectType {
		case globalkey.SingleMsg:
			err := l.svcCtx.SingleMsgModel.SoftDelete(nil, &database.SingleMsg{
				Id:               seq,
				Status:           globalkey.Del,
			})
			if err != nil {
				return nil, xerr.NewErrCodeMsg(xerr.USER_OPERATION_ERR, fmt.Sprintf("%+v删除成功，到%d删除失败, err:%s", success, seq, err.Error()))
			}
		case globalkey.GroupMsg:
			err := l.svcCtx.GroupMsgModel.SoftDelete(nil, &database.GroupMsg{
				Id:               seq,
				Status:           globalkey.Del,
			})
			if err != nil {
				return nil, xerr.NewErrCodeMsg(xerr.USER_OPERATION_ERR, fmt.Sprintf("%+v删除成功，到%d删除失败, err:%s", success, seq, err.Error()))
			}
		}
		success = append(success, seq)
	}
	return &chat.NullResp{}, nil
}
