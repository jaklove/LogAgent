package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"learn-log-agent/pkg/config"
	"log"
	"time"
)

var (
	EctdClient *clientv3.Client
)

type CollectEntry struct {
	Path string `json:"path"`
	Topic string `json:"topic"`
}

func Init()error{
	client, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.EtcdCfg.Address},
		DialTimeout: time.Second * 5,
	})
	if e != nil{
		fmt.Println("init etcd config fail")
		return e
	}
	EctdClient = client
	return nil
}

func GetConf()([]CollectEntry,error){
	var CollectList []CollectEntry
	ctx,canel := context.WithTimeout(context.Background(),time.Second*5)
	defer canel()

	getResponse, err := EctdClient.Get(ctx, config.EtcdCfg.CollectKey)
	if err != nil{
		log.Println("get conf from ectd by key:%s,err:%v",config.EtcdCfg.CollectKey,err)
		return nil,err
	}
	if len(getResponse.Kvs) == 0{
		log.Printf("get len:0 conf form etcd by key:%s",config.EtcdCfg.CollectKey)
		return nil,nil
	}
	fmt.Println("config.EtcdCfg.CollectKey:",config.EtcdCfg.CollectKey)
	ret := getResponse.Kvs[0]
	err = json.Unmarshal(ret.Value,&CollectList)
	if err != nil{
		fmt.Println("json Unmarshal failed,err:",err)
		return nil,err
	}
	return CollectList,nil
}