package config

import "github.com/go-ini/ini"

type KafkaConfig struct {
	KafkaAddress string `ini:"address"`
	Topic string `ini:"topic"`
}

type CollectConfig struct {
	LogFilePath string `ini:"log_file_path"`
}

type EtcdConfig struct {
	Address string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

//定义配置文件
var KfkaConfig = &KafkaConfig{}
var CollectCfg = &CollectConfig{}
var EtcdCfg = &EtcdConfig{}

//读取配置
func InitConfig() error {
	cfg, e := ini.Load("./conf/conf.ini")
	if e != nil{
		return e
	}
	cfg.Section("kafka").MapTo(KfkaConfig)
	cfg.Section("collect").MapTo(CollectCfg)
	cfg.Section("etcd").MapTo(EtcdCfg)
	return nil
}
