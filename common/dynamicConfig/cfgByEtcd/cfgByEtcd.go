package cfgByEtcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// etcd client put/get demo
// use etcd/clientv3

type ConfigByEtcd struct {
	host string
	cli  *clientv3.Client
}

func NewConfigByEtcd(host string) *ConfigByEtcd {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{host},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(fmt.Sprintf("connect to etcd failed, err:%v\n", err) )
	}
	return &ConfigByEtcd{
		host: host,
		cli: cli,
	}
}

func (ce ConfigByEtcd)Close() error {
	return ce.cli.Close()
}

func (ce ConfigByEtcd)Put() {
	//defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := ce.cli.Put(ctx, "lmh", "lmh")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}

}

func (ce ConfigByEtcd)Get()  {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := ce.cli.Get(ctx, "lmh")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}

func (ce ConfigByEtcd)Watch() {
	//cli, err := clientv3.New(clientv3.Config{
	//	Endpoints:   []string{"127.0.0.1:2379"},
	//	DialTimeout: 5 * time.Second,
	//})
	//if err != nil {
	//	fmt.Printf("connect to etcd failed, err:%v\n", err)
	//	return
	//}
	//fmt.Println("connect to etcd success")
	//defer cli.Close()
	// watch key:lmh change
	rch := ce.cli.Watch(context.Background(), "lmh") // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}