package serverInfo

import (
	"runtime"
	"strconv"
)

type SysInfo struct {
	NumCpu          string `json:"num_Cpu"`
	NumCpuUsage     string `json:"num_Cpu_Usage"`
	NumRam          string `json:"num_Ram"`
	NumRamUsage     string `json:"num_Ram_Usage"`
	NumDisk         string `json:"num_Disk"`
	NumDiskUsage    string `json:"num_Disk_Usage"`
	NumNetwork      string `json:"num_Network"`
	NumNetworkUsage string `json:"num_Network_Usage"`
}

func GetSysInfo() *SysInfo {
	return &SysInfo{
		NumCpu:          strconv.Itoa(runtime.NumCPU()),
		NumCpuUsage:     "",
		NumRam:          "",
		NumRamUsage:     "",
		NumDisk:         "",
		NumDiskUsage:    "",
		NumNetwork:      "",
		NumNetworkUsage: "",
	}
}
