package tcp

import (
	"context"
	"fmt"
	"time"

	//"go.uber.org/zap"
	"github.com/alberliu/gn"
	"github.com/zeromicro/go-zero/core/logx"

	"im-center/common/tool"
	"im-center/common/uniqueid"
	"im-center/service/connect/core/connect"
	"im-center/service/connect/core/connectManager"
	"im-center/service/connect/internal/config"
	"im-center/service/connect/internal/types"
)
type tcpHandler struct{
	logx.Logger
	cm *connectManager.ConnectManager
}

var (
	tcpH = tcpHandler{
		logx.WithContext(context.Background()),
		connectManager.GetCM(),
	}
	server *gn.Server
)

// 启动TCP服务器
func TCPServer(c *config.Config) {
	gn.SetLogger(nil)
	server, err := gn.NewServer(c.TCPListenAddr, &tcpH,
		gn.WithDecoder(gn.NewHeaderLenDecoder(2)),
		gn.WithEncoder(gn.NewHeaderLenEncoder(2, 1024)),
		gn.WithReadBufferLen(256),
		gn.WithTimeout(11*time.Minute),
		gn.WithAcceptGNum(10),
		gn.WithIOGNum(100),
	)
	if err != nil {
		fmt.Print("TCP服务器启动失败:", err)
		panic(err)
	}

	fmt.Printf("Starting tcp server at %s...\n", c.TCPListenAddr)
	server.Run()
}

//停止TCP服务器
func StopTCPServer() {
	server.Stop()
}

func (t *tcpHandler) OnConnect(c *gn.Conn) {
	// 鉴权+通过token获取userId、deviceId

	userId := uniqueid.GenUid()	// todo 真实userid
	deviceId := "api"
	if exist := t.cm.ExistConnect(&types.ConnectUid{
		UserId: userId,
		DeviceId: deviceId,
	}); exist {
		t.Errorf("userId:%d, deviceId:%s, 该用户设备的tpc长连接已存在", userId, deviceId)
		c.Close()
		return
	}
	c.SetData(&types.ConnectUid{
		UserId: userId,
		DeviceId: deviceId,
	})

	// 初始化连接数据
	conn := NewTcpConnect(context.Background(), tcpH.cm.Config, types.TcpConnect{
		Socket: c,
		Addr:        c.GetAddr(),
		UserId:      userId,
		DeviceId:    deviceId,
		CurrentTime: tool.GetUintNowTime(),
	})
	t.Infof("connect, fd:%d, addr:%s", c.GetFd(), c.GetAddr() )

	if t.cm.Config.ConnRWSeparate{
		go conn.Write()
	}

	// 用户连接事件
	t.cm.Register <- conn
}

func (t *tcpHandler) OnMessage(c *gn.Conn, message []byte) {
	connectUid := c.GetData().(*types.ConnectUid)
	if conn := t.cm.GetConnect(connectUid); conn != nil{
		// 处理程序
		connect.ProcessData(conn, message)
	} else {
		t.Errorf("connectUid:%+v, 连接不存在", connectUid)
	}
}

func (t *tcpHandler) OnClose(c *gn.Conn, err error) {
	connectUid := c.GetData().(*types.ConnectUid)
	if conn := t.cm.GetConnect(connectUid); conn != nil{
		t.cm.Unregister <- conn
	} else {
		t.Errorf("connectUid:%+v, 连接不存在", connectUid)
	}
}
