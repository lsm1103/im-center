//连接管理器的对外统一接口
package connectManager

import (
	"errors"
	"im-center/common/globalkey"
	"im-center/common/tool"
	"im-center/service/connect/core/connect"
	"im-center/service/connect/internal/types"
	"im-center/service/connect/internal/utils"
	"strings"
)

var (
	connectExistErr   = errors.New("该节点已存在该连接")
	connectAddErr     = errors.New("建立连接失败")
	connectNotFindErr = errors.New("找不到该连接")
)

func (cm *ConnectManager) SendOneMsg(req *types.SendOneMsgReq) error {
	var receiverDeviceIds []string
	if req.ReceiverDeviceId != "*"{
		receiverDeviceIds = strings.Split(req.ReceiverDeviceId, ",")
	} else {
		receiverDeviceIds = cm.GetUserConnectIdList(req.ReceiverId)
	}
	if len(receiverDeviceIds) == 0 {
		return connectNotFindErr
	}
	for _,id := range receiverDeviceIds {
		conn := cm.getConnect(&types.ConnectUid{
			UserId: req.ReceiverId,
			DeviceId: id,
		})
		if conn == nil {
			cm.Errorf("%s:%s 找不到该连接", req.ReceiverId, id)
		}
		msg := utils.BuildMsg(req)
		if msg == nil {
			cm.Errorf("%s:%s 格式化消息失败", req.ReceiverId, id)
		}
		conn.SendMsg(msg )
	}
	return nil
}

func (cm *ConnectManager) SendManyMsg(req *types.SendManyMsgReq) error {
	for _,id := range req.ReceiverId {
		err := cm.SendOneMsg(&types.SendOneMsgReq{
			SenderType:       req.SenderType,
			SenderId:         req.SenderId,
			SenderDeviceId:   req.SenderDeviceId,
			ReceiverId:       id,
			ReceiverDeviceId: req.ReceiverDeviceId,
			ParentId:         req.ParentId,
			SendTime:         req.SendTime,
			MsgType:          req.MsgType,
			MsgContent:       req.MsgContent,
		})
		if err != nil {
			cm.Errorf("给%s发送消息失败，err:%s", id, err.Error())
		}
	}

	return nil
}

// 主动关闭连接，仅供外步调用
func (cm *ConnectManager) OffConnect(req *types.ConnectUid) error {
	conn := cm.getConnect(req)
	if conn != nil {
		return connectNotFindErr
	}
	conn.ActiveClose()
	return nil
	//return cm.eventUnregister(conn)
}

// 查询连接是否存在, 可以实现查询用户是否在线功能
func (cm *ConnectManager) ExistConnect(req *types.ConnectUid) bool {
	if _,ok := cm.Connects.Load(globalkey.BuildKey(globalkey.ConnectIdFmt, req.UserId, req.DeviceId)); ok {
		return true
	}
	return false
}

// 获取连接详情
func (cm *ConnectManager) GetConnectInfo(req *types.ConnectUid) (resp *types.ConnectInfoResp, err error) {
	if obj := cm.getConnect(req); obj != nil {
		data := obj.GetData()
		return &types.ConnectInfoResp{
			UserId:         data.UserId,
			DeviceId:       data.DeviceId,
			ServerIp:       data.ServerIp,
			ConnectIp:      data.Addr,
			RegisterTime:   data.RegisterTime,
			HeartbeatTime:  data.HeartbeatTime,
		}, nil
	}
	return nil, connectNotFindErr
}

// 分页获取连接列表
func (cm *ConnectManager) GetConnectList(req *types.GetConnectListReq) (resp *types.ConnectListResp, err error) {
	tmp := []types.ConnectItem{}
	cm.Connects.Range(func(key, value interface{}) bool {
		conn := value.(connect.Connect)
		data := conn.GetData()
		tmp = append(tmp, types.ConnectItem{
			UserId:         data.UserId,
			DeviceId:       data.DeviceId,
			ServerIp:       data.ServerIp,
			ConnectIp:      data.Addr,
			RegisterTime:   data.RegisterTime,
			HeartbeatTime:  data.HeartbeatTime,
		})
		return true
	})

	return &types.ConnectListResp{
		ConnectList: tmp,
	}, nil
}

// 分页获取在线用户列表
func (cm *ConnectManager) GetOnlineUserList(req *types.GetOnlineUserListReq) (resp *types.OnlineUserListResp, err error) {
	cm.UserLock.Lock()
	var count uint64 = 0
	tmp := []types.OnlineUserItem{}
	for k,v := range cm.Users {
		if req.Offset < count && count < req.Offset+req.Limit {
			tmp = append(tmp, types.OnlineUserItem{
				UserId: k,
				OnlineDevices: v,
			})
		}
		count += 1
	}
	cm.UserLock.Unlock()
	return &types.OnlineUserListResp{
		OnlineUserList: tmp,
	}, nil
}

// 获取某个用户上线的设备列表
func (cm *ConnectManager) GetUserConnectList(req *types.GetUserConnectListReq) (resp *types.ConnectListResp, err error) {
	cm.UserLock.Lock()
	connectIds, ok := cm.Users[req.UserId]
	cm.UserLock.Unlock()
	if !ok{ return }

	tmp := []types.ConnectItem{}
	for _,connectId := range connectIds{
		conn, _ := cm.Connects.Load(connectId)
		if conn != nil{
			conn_ := conn.(connect.Connect)
			data := conn_.GetData()
			tmp = append(tmp, types.ConnectItem{
				UserId:         data.UserId,
				DeviceId:       data.DeviceId,
				ServerIp:       data.ServerIp,
				ConnectIp:      data.Addr,
				RegisterTime:   data.RegisterTime,
				HeartbeatTime:  data.HeartbeatTime,
			})
		}
	}

	return &types.ConnectListResp{
		ConnectList: tmp,
	}, nil
}

// 获取某个用户上线的设备列表
func (cm *ConnectManager) GetUserConnectIdList(userId string) (resp []string) {
	cm.UserLock.Lock()
	resp, ok := cm.Users[userId]
	cm.UserLock.Unlock()
	if !ok{ return nil }
	return
}

// 获取业务信息
func (cm *ConnectManager) GetBusinessInfo() *types.BusinessInfo {
	return &types.BusinessInfo{
		ConnectLen:        cm.ConnectNum,
		UserLen:           int64(len(cm.Users)),
		PendRegisterLen:   int64(len(cm.Register)),
		PendUnregisterLen: int64(len(cm.Unregister)),
	}
}

// 获取连接
func (cm *ConnectManager) GetConnect(req *types.ConnectUid) connect.Connect {
	return cm.getConnect(req)
}

// 存储连接
func (cm *ConnectManager) saveConnect(req connect.Connect) error {
	// 检查是否已经存在
	cData := req.GetData()
	connectId := globalkey.BuildKey(globalkey.ConnectIdFmt, cData.UserId, cData.DeviceId)
	if _, loaded := cm.Connects.LoadOrStore(connectId, req); loaded {
		return connectExistErr
	}

	cm.UserLock.Lock()
	cm.ConnectNum += 1
	_,ok := cm.Users[cData.UserId]
	if ok{
		cm.Users[cData.UserId] = append(cm.Users[cData.UserId], connectId)
	} else {
		cm.Users[cData.UserId] = []string{connectId,}
	}
	cm.UserLock.Unlock()
	return nil
}
// 删除连接
func (cm *ConnectManager) deleteConnect(req *types.ConnectUid) {
	connectId := globalkey.BuildKey(globalkey.ConnectIdFmt, req.UserId, req.DeviceId)
	cm.Connects.Delete(connectId)
	if cm.ConnectNum > 0{
		cm.ConnectNum -= 1
	}

	cm.UserLock.Lock()
	cm.Users[req.UserId] = tool.DelStrSlice(cm.Users[req.UserId], connectId)
	if len(cm.Users[req.UserId]) == 0 {
		delete(cm.Users, req.UserId)
	}
	cm.UserLock.Unlock()
}
// 获取连接
func (cm *ConnectManager) getConnect(req *types.ConnectUid) connect.Connect {
	if val,ok := cm.Connects.Load(globalkey.BuildKey(globalkey.ConnectIdFmt, req.UserId, req.DeviceId)); ok {
		return val.(connect.Connect)
	}
	return nil
}
// 定时清理超时连接
func (cm *ConnectManager) ClearTimeoutConnections() {
	cm.Infof("所有连接数：%+v, UsersNum：%+v, 未处理上线数：%+v, 未处理下线数：%+v", cm.ConnectNum, len(cm.Users), len(cm.Register), len(cm.Unregister))
	currentTime := tool.GetUintNowTime()
	cm.Connects.Range(func(connectId, conn interface{}) bool {
		conn_ := conn.(connect.Connect)
		if conn_.IsHeartbeatTimeout(currentTime) {
			cm.Info("心跳时间超时 关闭连接", conn_.GetData())
			// 下线
			conn_.ActiveClose()
		}
		return true
	})
}



// 循环本地该用户的所有连接
//func (cm *ConnectManager) loopUserAllConnect(userId string, f func(conn *connect.Connect)) (count int) {
//	cm.UserLock.Lock()
//	connects, ok := cm.Users[userId]
//	cm.UserLock.Unlock()
//	if !ok{ return }
//	count = len(connects)
//	for _,connectId := range connects{
//		conn, _ := cm.Connects.Load(connectId)
//		if conn != nil{
//			f(conn.(*connect.Connect))
//		}
//	}
//	return
//}