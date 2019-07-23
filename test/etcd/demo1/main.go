package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {

	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		putResp *clientv3.PutResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job11", "echo helloq", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(putResp.Header.Revision)
		fmt.Println(string(putResp.PrevKv.Value))
	}

}
