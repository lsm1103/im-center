package logic

import (
	"context"
	"im-center/service/connect/internal/types"
	"strconv"

	"im-center/service/connect/rpc/connect"
	"im-center/service/connect/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConnectListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserConnectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConnectListLogic {
	return &GetUserConnectListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//获取用户连接列表,分页
func (l *GetUserConnectListLogic) GetUserConnectList(in *connect.GetUserConnectListReq) (*connect.UserConnectListResp, error) {
	list, err := l.svcCtx.CS.GetUserConnectList(&types.GetUserConnectListReq{
		UserId: in.UserId,
	})
	if err != nil {
		return nil, err
	}

	resp := &connect.UserConnectListResp{
		UserConnectList: []*connect.ConnectInfo{},
	}
	for _,item := range list.ConnectList {
		resp.UserConnectList = append(resp.UserConnectList, &connect.ConnectInfo{
			UserId:               item.UserId,
			DeviceId:             item.DeviceId,
			ServerIp:             item.ServerIp,
			ConnectIp:            item.ConnectIp,
			RegisterTime:   	  strconv.FormatUint(item.RegisterTime, 10),
			HeartbeatTime:  	  strconv.FormatUint(item.HeartbeatTime, 10),
			UnRegisterTime: 	  strconv.FormatUint(item.UnRegisterTime, 10),
			DeviceInfo:           item.DeviceInfo,
		})
	}
	return resp, nil
}
