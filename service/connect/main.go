package main

import (
	"flag"
	"fmt"
	"im-center/service/connect/internal/timedTask"
	"runtime"

	"github.com/zeromicro/go-zero/core/conf"
	"im-center/service/connect/api"
	"im-center/service/connect/core"
	"im-center/service/connect/internal/config"
	"im-center/service/connect/internal/svc"
	"im-center/service/connect/rpc"
)

/*
1、im服务不登入，通过注册携带的token进行鉴权
2、长链接服务和发信息、用户、群组服务分离；长链接服务只支持注册上线、下线、消息下发
3、消息下发通过grpc
4、多节点通过grpc做消息转发
 */

var apiConfigFile = flag.String("apiF", "etc/connect_dev.yaml", "the config file")
var rpcConfigFile = flag.String("rpcF", "etc/rpc_dev.yaml", "the config file")

func main() {
	flag.Parse()
	fmt.Printf("numGoroutine: %d\n", runtime.NumGoroutine())
	var c config.Config
	conf.MustLoad(*apiConfigFile, &c)
	fmt.Printf("numGoroutine: %d\n", runtime.NumGoroutine())

	// 初始化服务上下文
	ctx := svc.NewServiceContext(c)

	// 初始化连接层服务
	cs := core.NewConnectServer(&c, ctx.Cache)
	cs.Start()
	defer cs.Stop()
	fmt.Printf("numGoroutine: %d\n", runtime.NumGoroutine())

	ctx.Cs = cs

	apiServer := api.Run(c, ctx)
	defer apiServer.Stop()

	fmt.Printf("numGoroutine: %d\n", runtime.NumGoroutine())
	go rpc.Run(*rpcConfigFile, cs)

	// 启动定时任务
	timedTask.Run(ctx)
	fmt.Printf("numGoroutine: %d\n", runtime.NumGoroutine())

	apiServer.Start()
}