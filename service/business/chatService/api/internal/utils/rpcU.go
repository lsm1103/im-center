package utils

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"im-center/common/timedTask"
	"im-center/service/business/chatService/api/internal/types"
	"im-center/service/connect/rpc/connectclient"
	"im-center/service/model/cache"
	"time"
)

type (
	RpcU struct {
		logx.Logger
		Cache *cache.RedisCache
		NodeList *types.ServerInfoResp
	}
)

func NewRpcU(ch *cache.RedisCache) *RpcU {
	r := &RpcU{
		Logger: logx.WithContext(context.Background()),
		Cache: ch,
	}
	go timedTask.Timer(3*time.Second, 30*time.Second, upNodeList, r, nil, nil)
	return r
}

func (r *RpcU) GetNode() connectclient.Connect {
	rpcIp := r.NodeList.Server[0].ServerId
	return connectclient.NewConnect(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{rpcIp},
		NonBlock: true,
	}))
}

func (r *RpcU) GetNodeList() (resp *types.ServerInfoResp, err error) {
	nodeList := r.Cache.GetNodeList()
	var resp_ []types.ServerItem
	for _,item := range nodeList {
		tmp := &types.ServerItem{}
		err = json.Unmarshal([]byte(item), tmp)
		if err != nil {
			r.Errorf("获取分布式节点列表：%+v,%+v", item, err)
			return nil, err
		}
		resp_ = append(resp_, *tmp)
	}
	r.Errorf("获取分布式节点列表：%+v", resp_)
	return &types.ServerInfoResp{
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