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

type GroupBatchDelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupBatchDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupBatchDelLogic {
	return &GroupBatchDelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupBatchDelLogic) GroupBatchDel(in *chat.GroupBatchDelReq) (*chat.NullResp, error) {
	success := []int64{}
	for _,id := range in.GroupIds {
		err := l.svcCtx.GroupModel.SoftDelete(nil, &database.Group{
			Id:     id,
			Status: globalkey.Del,
		})
		if err != nil {
			return nil, xerr.NewErrCodeMsg(xerr.USER_OPERATION_ERR, fmt.Sprintf("%+v删除成功，到%d删除失败, err:%s", success, id, err.Error()))
		}
		success = append(success, id)
	}
	return &chat.NullResp{}, nil
}
