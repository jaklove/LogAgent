package tailfile

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	config2 "learn-log-agent/pkg/config"
	"learn-log-agent/pkg/kafka"
	"log"
	"strings"
	"time"
)

var (
	TailObj *tail.Tail
)




func Init()error{
	var err error
	config := tail.Config{
		ReOpen:true,
		Follow:true,
		Location:&tail.SeekInfo{Offset:0,Whence:2},
		MustExist:false,
		Poll:true,
	}

	//打开文件开始读取数据
	TailObj,err = tail.TailFile(config2.CollectCfg.LogFilePath,config)
	if err != nil{
		fmt.Printf("tail %s failed,err:%v\n",config2.CollectCfg.LogFilePath,err)
		return err
	}

	//后台读取日志文件发送到kafka中
	go run()
	return nil
}

func run()  {
	for  {
		//循环读取文件
		line,ok := <- TailObj.Lines
		if !ok{
			log.Println("tail file close reopen,filename:",config2.CollectCfg.LogFilePath)
			time.Sleep(1*time.Second)
			continue
		}
		fmt.Println("line.Text:",line.Text)
		if len(strings.Trim(line.Text,"\r")) == 0{
			continue
		}

		//封装消息
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web-log"
		msg.Value = sarama.StringEncoder(line.Text)
		//发送到chan中
		kafka.MsgChan <- msg
	}
}
