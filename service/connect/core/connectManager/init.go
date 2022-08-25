package connectManager

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"im-center/service/connect/core/connect"
	"im-center/service/connect/core/init/distributed"
	"im-center/service/connect/internal/config"
	"sync"
)

// 连接管理
type ConnectManager struct {
	logx.Logger
	Config      *config.Config
	Distributed *distributed.Distributed

	Connects 	sync.Map		   // 全部的连接 { ConnectId:*Connect }; ConnectId: 节点ip端口+用户id+设备id
	ConnectNum	int64
	Users       map[string][]string // 全部的用户设备连接 { userid:[ConnectId, ConnectId, ] };
	UserLock    sync.RWMutex       // 读写锁
	Register    chan connect.Connect       // 连接连接处理
	Unregister  chan connect.Connect       // 断开连接处理程序
}

var connectManager = &ConnectManager{
	Logger: 	logx.WithContext(context.Background()),
	Users:      make(map[string][]string),
	Register:   make(chan connect.Connect, 1000),
	Unregister: make(chan connect.Connect, 1000),
	//Broadcast:  make(chan []byte, 1000),
}

func GetCM() *ConnectManager {
	return connectManager
}

// 启动
func (cm *ConnectManager) Start() {
	for {
		select {
		case conn := <-cm.Register:
			// 建立连接事件
			cm.eventRegister(conn)
		case conn := <-cm.Unregister:
			// 断开连接事件
			err := cm.eventUnregister(conn)
			if err != nil {
				cm.Errorf("eventUnregister error: %s", err.Error())
			}
		}
	}
}

// 关闭
func (cm *ConnectManager) Stop() {
	close(cm.Register)
	close(cm.Unregister)
}

func (cm *ConnectManager)SetRelyOn(Config *config.Config, Distributed *distributed.Distributed) {
	cm.Config = Config
	cm.Distributed = Distributed
}
