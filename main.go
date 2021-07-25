package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"learn-log-agent/pkg/config"
	"learn-log-agent/pkg/etcd"
	"learn-log-agent/pkg/kafka"
	"learn-log-agent/pkg/tailfile"
	"log"
	"time"
)

func init() {
	var err error

	//初始化配置
	err = config.InitConfig()
	if err != nil {
		log.Fatalf("loading config file err:%s", err.Error())
	}

	//初始化kafka连接
	err = kafka.InitKafkaConfig()
	if err != nil {
		log.Fatalf("loading kafka config file err:%s", err.Error())
	}

	//加载etcd配置文件
	err = etcd.Init()
	if err != nil{
		log.Fatalf("loading etcd.init file err:%s", err.Error())
	}
	collectEntries, err := etcd.GetConf()
	if err != nil{
		log.Fatalf("etcd.GetConf  err:%s", err.Error())
	}
	fmt.Println("collectEntries:",collectEntries)

	//加载tail
	err = tailfile.Init()
	if err != nil {
		log.Fatalf("loading tailfile config file err:%s", err.Error())
	}
}

func main() {
	var err error
	client, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.111.99:2379"},
		DialTimeout: time.Second * 5,
	})

	if e != nil{
		fmt.Println("connect to etcd failed err:",e)
	}

	defer client.Close()

	//执行put命令
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	configstr := `[{"path":"E:\\logs\\log.log","topic":"log"},{"path":"E:\\logs\\log.log1","topic":"log1"}]`
	_,err = client.Put(ctx,"log_zrj",configstr)
	if err != nil{
		panic(err)
	}

	_, err = client.Put(ctx, "zrj-etcd", "helloword zhourenjie")
	if err != nil{
		fmt.Println("put etcd key faild err:",err)
	}

	//取值
	gctx, gcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer gcancel()

	getResponse, err := client.Get(gctx, "zrj-etcd")
	if err != nil{
		fmt.Println("get etcd key err",err)
		return
	}

	for _,ev := range getResponse.Kvs{
		fmt.Printf("key: %s  value:%s",ev.Key,ev.Value)
	}

	//监听watch
	watchCh := client.Watch(context.Background(), "s4")
	for resp := range watchCh{
		for _,ev := range resp.Events{
			fmt.Printf("type :%s, key :%s value:%s\n",ev.Type,ev.Kv,ev.Kv.Value )
		}
	}

	fmt.Println("success")
	for {
		time.Sleep(time.Second * 1)
	}
}
