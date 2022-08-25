package tool

import (
	"encoding/json"
	"runtime"
)

// 获取本系统的协程数、gc数和cpu、内存、硬盘使用情况
func GetSysInfo() (systemInfo map[string]interface{}) {
	systemInfo = make(map[string]interface{})

	systemInfo["goroutineNum"] = runtime.NumGoroutine() // 协程数
	systemInfo["goroutineInfo"] = GetStackInfo() // 协程数
	//systemInfo["gcNum"] = runtime.NumGc()                // gc数
	systemInfo["cpuPercent"] = runtime.NumCPU()          // cpu使用率
	//systemInfo["memory"] = runtime.MemStats().Alloc      // 内存使用量
	//systemInfo["disk"] = tool.GetDiskUsage()             // 硬盘使用量
	return
}

func GetStackInfo() string {
	buf := make([]byte, 64 * 1024)
	runtime.Stack(buf, true)
	return string(buf)
}

func getSysStats() (string, error) {
	sysStats := runtime.MemStats{}
	runtime.ReadMemStats(&sysStats)
	r, err := json.Marshal(sysStats)
	return string(r), err
}