package serverInfo

import (
	"fmt"
	"im-center/common/tool"
	"runtime"
	"strconv"
)

type ServiceInfo struct {
	NumGoroutine     string `json:"num_goroutine"`
	AllocMemory      string `json:"alloc_memory"`
	TotalAllocMemory string `json:"totalAlloc_memory"`
	SysMemory        string `json:"sys_memory"`
	NumGC            string `json:"num_GC"`
}

func GetServiceInfo() *ServiceInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &ServiceInfo{
		NumGoroutine:     strconv.Itoa(runtime.NumGoroutine()),
		AllocMemory:      fmt.Sprintf("%dMB",tool.BToMb(m.Alloc)),
		TotalAllocMemory: fmt.Sprintf("%dMB",tool.BToMb(m.TotalAlloc)),
		SysMemory:        fmt.Sprintf("%dMB",tool.BToMb(m.Sys)),
		NumGC: 			  strconv.Itoa(int(m.NumGC)),
	}
}
