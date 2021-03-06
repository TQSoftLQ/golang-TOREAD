/*
etcdctl存的数据，使用clientv3API也访问不到。etcdctl默认是使用v2 版本，所有在执行命令前需要设置环境变量，export ETCDCTL_API=3
v2和v3 不在一个存储区，两者之间是数据隔离的

$env:ETCDCTL_API=3


*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("连接ectd服务器集群失败")
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	respPut, err := cli.Put(ctx, "test_key", "etcd存储的值") //设置
	cancel()

	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx 被其他携程取消 : %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx 已超时: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("连接集群节点, which are not etcd servers: %v", err)
		}
	}

	respPut = respPut

	// fmt.Println(respPut.PrevKv)

	// use the response

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	respGet, err := cli.Get(ctx, "test_key") //获取etcd中内容，单次查询
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx 被其他携程取消 : %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx 已超时: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("连接集群节点, which are not etcd servers: %v", err)
		}
	}

	for k, v := range respGet.Kvs {
		fmt.Printf("index:%v key:%v value:%v\n", k, string(v.Key), string(v.Value))
	}
}
