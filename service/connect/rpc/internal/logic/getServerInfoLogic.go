package logic

import (
	"context"
	"im-center/service/connect/internal/types"
	"im-center/service/connect/rpc/connect"
	"im-center/service/connect/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetServerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerInfoLogic {
	return &GetServerInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetServerInfoLogic) GetServerInfo(in *connect.ServerInfoReq) (*connect.ServerInfoResp, error) {
	info, err := l.svcCtx.CS.GetServiceInfo(&types.GetServerInfoReq{
		ServerId: in.ServerId,
	})
	if err != nil {
		return nil, err
	}

	resp := []*connect.Server{}
	for _,item := range info.Server {
		resp = append(resp, &connect.Server{
			ServerId:             item.ServerId,
			ServerInfo:           &connect.ServerInfo{
				NumGoroutine:         item.ServerInfo.NumGoroutine,
				AllocMemory:          item.ServerInfo.AllocMemory,
				TotalAllocMemory:     item.ServerInfo.TotalAllocMemory,
				SysMemory:            item.ServerInfo.SysMemory,
				Num_GC:               item.ServerInfo.NumGC,
			},
			BusinessInfo:         &connect.BusinessInfo{
				ConnectLen:        item.BusinessInfo.ConnectLen,
				UserLen:           item.BusinessInfo.UserLen,
				PendRegisterLen:   item.BusinessInfo.PendRegisterLen,
				PendUnregisterLen: item.BusinessInfo.PendUnregisterLen,
			},
			SysInfo:              &connect.SysInfo{
				Num_CPU:              item.SysInfo.NumCpu,
				Num_CPU_USAGE:        item.SysInfo.NumCpuUsage,
				Num_RAM:              item.SysInfo.NumRam,
				Num_RAM_USAGE:        item.SysInfo.NumRamUsage,
				Num_DISK:             item.SysInfo.NumDisk,
				Num_DISK_USAGE:       item.SysInfo.NumDiskUsage,
				Num_NETWORK:          item.SysInfo.NumNetwork,
				Num_NETWORK_USAGE:    item.SysInfo.NumNetworkUsage,
			},
		})
	}

	return &connect.ServerInfoResp{
		Server: resp,
	}, nil
}
