package timedTask

import (
	"fmt"
	"im-center/common/serverInfo"
	"im-center/service/connect/core/connectManager"
	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"
	"runtime/debug"
)

// 更新本服务信息
func upServiceInfo(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("定时任务-更新本服务信息...")
	result = true

	serviceInfo := serverInfo.GetServiceInfo()
	sysInfo := serverInfo.GetSysInfo()
	svc_ := param.(*svc.ServiceContext)
	svc_.Cs.Distributed.NodeRegister(&types.ServerItem{
		ServerId: 	  svc_.Config.GetServerIp(),
		ServerInfo:   types.ServerInfo{
			NumGoroutine:     serviceInfo.NumGoroutine,
			AllocMemory:      serviceInfo.AllocMemory,
			TotalAllocMemory: serviceInfo.TotalAllocMemory,
			SysMemory:        serviceInfo.SysMemory,
			NumGC: 			  serviceInfo.NumGC,
		},
		SysInfo:      types.SysInfo{
			NumCpu:          sysInfo.NumCpu,
			NumCpuUsage:     sysInfo.NumCpuUsage,
			NumRam:          sysInfo.NumRam,
			NumRamUsage:     sysInfo.NumRamUsage,
			NumDisk:         sysInfo.NumDisk,
			NumDiskUsage:    sysInfo.NumDiskUsage,
			NumNetwork:      sysInfo.NumNetwork,
			NumNetworkUsage: sysInfo.NumNetworkUsage,
		},
		BusinessInfo: *connectManager.GetCM().GetBusinessInfo(),
	})
	return
}
