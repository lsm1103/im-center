package logic

import (
	"context"
	"fmt"
	"im-center/common/xerr"
	"im-center/service/business/chatService/rpc/chat"
	"im-center/service/business/chatService/rpc/internal/svc"
	"im-center/service/model/database"

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
	success := []int64{}
	for _,id := range in.FriendIds{
		err := l.svcCtx.FriendModel.SoftDelete(nil, &database.Friend{
			Id:     id,
			Status: in.OperationType,
		})
		if err != nil {
			return nil, xerr.NewErrCodeMsg(xerr.USER_OPERATION_ERR, fmt.Sprintf("%+v修改成功，到%d修改失败, err:%s", success, id, err.Error()))
		}
		success = append(success, id)
	}
	return &chat.NullResp{}, nil
}
