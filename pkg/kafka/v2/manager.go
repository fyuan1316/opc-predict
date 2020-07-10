package v2

import (
	"fmt"
	"opcdata-predict/cmd/option"
	"opcdata-predict/pkg/server"
)

type ConsumerManager struct {
	Consumer  *KafkaConsumer
	CommandCh chan []byte
}

func NewConsumerManager(opts option.Options) *ConsumerManager {
	m := &ConsumerManager{}
	m.Consumer = NewKafkaConsumer(opts)
	m.CommandCh = make(chan []byte)
	return m
}

func (cm *ConsumerManager) Process(dataCh chan<- []byte) {

	for {
		select {
		case bMsg := <-cm.CommandCh:
			if string(bMsg) == server.ControlCommand.Start {
				fmt.Println("recv start consumer command")
				//if cm.Consumer.Client == nil {
				//	if err := cm.Consumer.CreateClient(); err != nil {
				//		fmt.Printf("kafka create client err,%v\n", err)
				//	} else {
				if err := cm.Consumer.Run(dataCh); err != nil {
					fmt.Printf("kafka consume err,%v\n", err)
					panic(err)
				}
				//}
				//}
			}
			if string(bMsg) == server.ControlCommand.Stop {
				fmt.Println("recv stop consumer command")
				cm.Consumer.Stop()
			}
		}
	}

}
