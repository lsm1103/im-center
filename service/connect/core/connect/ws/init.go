package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"im-center/common/auth"
	"im-center/common/tool"
	"im-center/service/connect/core/connectManager"
	"im-center/service/connect/internal/svc"
	"im-center/service/connect/internal/types"
	"net/http"
	"strings"
)

func WsServer(w http.ResponseWriter, req *http.Request, svcCtx *svc.ServiceContext) {
	fmt.Println("升级协议", "ua:", req.Header["User-Agent"], "referer:", req.Header["Referer"])
	protocolStr := req.Header.Get("Sec-WebSocket-Protocol")
	if protocolStr == "" {
		http.Error(w, "连接失败, 未携带token", http.StatusUnauthorized)
		return
	}
	protocols := strings.Split(protocolStr, ",")
	tokenData, err := auth.ParseJwtToken(protocols[0], svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		http.Error(w, "连接失败, token不合法", http.StatusUnauthorized)
		return
	}
	// todo 鉴权+通过token获取userId、deviceId
	// 升级协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: protocols,
	}).Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, "连接失败", http.StatusMethodNotAllowed)
		return
	}

	cm := connectManager.GetCM()
	userId := tokenData["userId"]
	deviceId := tokenData["deviceId"]
	if exist := cm.ExistConnect(&types.ConnectUid{
		UserId: userId,
		DeviceId: deviceId,
	}); exist {
		http.Error(w, "该用户设备的长连接已存在", http.StatusBadRequest)
		return
	}

	connect := NewWsConnect(req.Context(), cm.Config, types.WsConnect{
		Socket:      conn,
		Addr:        conn.RemoteAddr().String(),
		UserId:      userId,
		DeviceId:    deviceId,
		CurrentTime: tool.GetUintNowTime(),
	})
	fmt.Printf("webSocket 建立连接:%s, UserId:%s, DeviceId:%s\n", conn.RemoteAddr().String(), userId, deviceId)

	// 连接是否读写分离协程
	if cm.Config.ConnRWSeparate{
		go connect.Read()
		go connect.Write()
	} else {
		go connect.Read()
	}

	// 用户连接事件
	cm.Register <- connect
}
