package rpc

import (
	"flag"
	"fmt"
	"im-center/service/connect/rpc/connect"

	"im-center/service/connect/core"
	"im-center/service/connect/rpc/internal/config"
	"im-center/service/connect/rpc/internal/server"
	"im-center/service/connect/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func Run(configFile string, cs *core.ConnectServer) {
	flag.Parse()

	var c config.Config
	conf.MustLoad(configFile, &c)
	ctx := svc.NewServiceContext(c, cs)
	srv := server.NewConnectServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		connect.RegisterConnectServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
