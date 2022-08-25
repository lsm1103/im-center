package connectManager

import (
	"im-center/common/globalkey"
	"im-center/service/connect/core/connect"
	"im-center/service/connect/internal/types"
)

// 建立连接事件
func (cm *ConnectManager) eventRegister(connect connect.Connect) {
	cm.Infof("建立连接 addr:%s", connect.GetData().Addr)
	//增加连接到连接管理器
	err := cm.saveConnect(connect)
	if err == connectExistErr {
		cm.Errorf("err:%+v",err)
		return
	}else if err != nil {
		cm.Errorf("err:%+v",connectAddErr)
		return
	}

	cm.Distributed.EventHandler(&types.DistributedEvent{
		EventType: globalkey.EventTypeOnline,
		ConnectData:      connect.GetData(),
	})

	// 发消息包到一个地方，通知该用户注册上线（可以作为一个关于用户/群组的配置，设置用户上线后是否需要通知、通知的范围）
}

// 断开连接事件
func (cm *ConnectManager) eventUnregister(connect connect.Connect) error {
	data := connect.GetData()
	cm.Infof("断开连接事件 UserId:%s,DeviceId:%s", data.UserId, data.DeviceId)
	// 关闭该连接的消息管道，来关闭为该连接开启的协程
	connect.PassiveClose()
	//从连接管理器删除连接
	cm.deleteConnect(&types.ConnectUid{
		UserId:   data.UserId,
		DeviceId: data.DeviceId,
	})

	cm.Distributed.EventHandler(&types.DistributedEvent{
		EventType: globalkey.EventTypeOffline,
		ConnectData:      data,
	})

	// 连接下线通知,发消息包到一个地方，通知该用户下线（可以作为一个关于用户/群组的配置，设置用户下线后是否需要通知、通知的范围）

	return nil
}

// 续期事件
func (cm *ConnectManager) EventRenewal(data *types.ConnectData) {
	//把连接从公共缓存集群中续期
	cm.Infof("续期事件 UserId:%s,DeviceId:%s", data.UserId, data.DeviceId)
	cm.Distributed.EventHandler(&types.DistributedEvent{
		EventType: globalkey.EventTypeRenewal,
		ConnectData:      data,
	})
}
