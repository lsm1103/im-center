package nodeHandler

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"im-center/common/timedTask"
	"im-center/service/connect/rpc/connectclient"
	"im-center/service/model/cache"
	"time"
)

type (
	RpcU struct {
		logx.Logger
		Cache    *cache.RedisCache
		NodeList *ServerInfoResp
	}

	ServerItem struct {
		ServerId     string       `json:"server_id"`     //服务器id
		ServerInfo   ServerInfo   `json:"server_info"`   //服务器信息
		BusinessInfo BusinessInfo `json:"business_info"` //业务信息
		SysInfo      SysInfo      `json:"sys_info"`      //系统信息
	}

	BusinessInfo struct {
		ConnectLen        string `json:"connect_len"`        //连接数
		UserLen           string `json:"user_len"`           //登录用户数
		PendRegisterLen   string `json:"pendRegister_len"`   //未处理连接事件数
		PendUnregisterLen string `json:"pendUnregister_len"` //未处理退出登录事件数
	}

	SysInfo struct {
		NumCpu          string `json:"num_Cpu"`           // Cpu数量
		NumCpuUsage     string `json:"num_Cpu_Usage"`     // Cpu使用率
		NumRam          string `json:"num_Ram"`           // 内存大小
		NumRamUsage     string `json:"num_Ram_Usage"`     // 内存使用率
		NumDisk         string `json:"num_Disk"`          // 磁盘大小
		NumDiskUsage    string `json:"num_Disk_Usage"`    // 磁盘使用率
		NumNetwork      string `json:"num_Network"`       // 网络流量
		NumNetworkUsage string `json:"num_Network_Usage"` // 网络流量使用率
	}

	ServerInfo struct {
		NumGoroutine     string `json:"num_goroutine"`     // goroutine数量
		AllocMemory      string `json:"alloc_memory"`      // 分配内存
		TotalAllocMemory string `json:"totalAlloc_memory"` // 分配内存总量
		SysMemory        string `json:"sys_memory"`        // 系统内存
		NumGC            string `json:"num_GC"`            // gc数量
		CPU              string `json:"CPU"`               // cpu使用率
	}

	ServerInfoResp struct {
		Server []ServerItem `json:"server"` //服务器列表
	}
)

func NewRpcU(ch *cache.RedisCache) *RpcU {
	r := &RpcU{
		Logger: logx.WithContext(context.Background()),
		Cache:  ch,
	}
	go timedTask.Timer(3*time.Second, 30*time.Second, upNodeList, r, nil, nil)
	return r
}

func (r *RpcU) GetNode() connectclient.Connect {
	rpcIp := r.NodeList.Server[0].ServerId
	return connectclient.NewConnect(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{rpcIp},
		NonBlock:  true,
	}))
}

func (r *RpcU) GetNodeList() (resp *ServerInfoResp, err error) {
	nodeList := r.Cache.GetNodeList()
	var resp_ []ServerItem
	for _, item := range nodeList {
		tmp := &ServerItem{}
		err = json.Unmarshal([]byte(item), tmp)
		if err != nil {
			r.Errorf("获取分布式节点列表：%+v,%+v", item, err)
			return nil, err
		}
		resp_ = append(resp_, *tmp)
	}
	r.Errorf("获取分布式节点列表：%+v", resp_)
	return &ServerInfoResp{
		Server: resp_,
	}, nil
}

func upNodeList(param interface{}) bool {
	r := param.(*RpcU)
	nodeList, err := r.GetNodeList()
	if err != nil {
		return false
	}
	r.NodeList = nodeList
	return true
}
