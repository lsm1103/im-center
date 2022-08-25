package tool

import (
	"github.com/zeromicro/go-zero/zrpc"
	"im-center/service/connect/rpc/connectclient"
)

func GetImCenterRpc(ipAddr string) connectclient.Connect {
	return connectclient.NewConnect(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{ipAddr, },
		NonBlock: true,
	}))
}
