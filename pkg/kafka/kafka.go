package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"learn-log-agent/pkg/config"
	"log"
)

var (
	KafkaClient sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

func InitKafkaConfig() error {
	var err error
	newConfig := sarama.NewConfig()
	newConfig.Producer.RequiredAcks = sarama.WaitForAll //leader、follow全部确认
	newConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	newConfig.Producer.Return.Successes = true

	//连接kafka
	KafkaClient, err = sarama.NewSyncProducer([]string{config.KfkaConfig.KafkaAddress}, newConfig)
	if err != nil{
		return err
	}
	fmt.Println("kafka success!")

	MsgChan = make(chan *sarama.ProducerMessage,10000)

	//发送消息到kafka
	go sendMsg()

	return nil
}

func sendMsg()  {
	for  {
		select {
		case msg := <- MsgChan :
			pid, offset, err := KafkaClient.SendMessage(msg)
			if err != nil{
				log.Printf("send log msg faild,err:",err)
				continue
			}
			log.Printf("send msg to kafka success pid:%v offset:%v",pid,offset)
		}
		
	}
}
