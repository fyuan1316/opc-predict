package v2

import (
	"fmt"
	"github.com/Shopify/sarama"
	"opcdata-predict/cmd/option"
	"opcdata-predict/pkg/predict"
	"opcdata-predict/pkg/scopelog"
	"strings"
	"sync"
	"time"
)

type KafkaConsumer struct {
	Config    *sarama.Config
	Brokers   string
	Topic     string
	Partition int
	Client    sarama.Client
	Offset    *int64
	StopCh    chan struct{}
	isReplay  bool
	// cli args
	Options option.Options
}

func NewKafkaConsumer() *KafkaConsumer {
	c := &KafkaConsumer{}
	c.Brokers = "localhost:9092"
	c.Topic = "test"
	c.Partition = 0
	c.Config = sarama.NewConfig()
	c.StopCh = make(chan struct{})
	return c
}

var (
	doOnce        sync.Once
	scopeConsumer = "Consumer"
)

func (c *KafkaConsumer) CreateClient() error {
	var client sarama.Client
	var err error
	doOnce.Do(func() {
		client, err = sarama.NewClient(strings.Split(c.Brokers, ","), c.Config)
		//p := sarama.OffsetOldest
		p := int64(565)
		c.Offset = &p
	})
	if err != nil {
		return err
	}
	if c.Client == nil {
		c.Client = client
	}
	return nil
}

func (c *KafkaConsumer) Stop() {
	close(c.StopCh)
}

func (c *KafkaConsumer) Run(dataCh chan<- []byte) error {
	var client sarama.Client
	var err error
	//doOnce.Do(func() {
		client, err = sarama.NewClient(strings.Split(c.Brokers, ","), c.Config)
		//p := sarama.OffsetOldest
		p := int64(565)
		c.Offset = &p
	//})

	c.Client = client
	//c.Options = opts
	consumer, err := sarama.NewConsumerFromClient(c.Client)
	if err != nil {
		return err
	}
	//if c.StopCh == nil {
	c.StopCh = make(chan struct{})
	//}
	go func() {
		for {
			if replay := c.consume(consumer, dataCh); !replay {
				break
			}
			fmt.Println("replay")
		}
	}()
	return nil
}

//
//func (c *KafkaConsumer) StartConsume(wsConn *websocket.Conn, opts option.Options) error {
//	c.Options = opts
//	consumer, err := sarama.NewConsumerFromClient(c.Client)
//	if err != nil {
//		return err
//	}
//	//if c.StopCh == nil {
//	c.StopCh = make(chan struct{})
//	//}
//	go func() {
//		for {
//			if replay := c.consume(consumer, wsConn); !replay {
//				break
//			}
//			fmt.Println("replay")
//		}
//	}()
//	return nil
//}

func (c *KafkaConsumer) getOffset() int64 {
	if c.isReplay {
		return sarama.OffsetOldest
	} else {
		return *c.Offset
	}
}

func (c *KafkaConsumer) consume(consumer sarama.Consumer, dataCh chan<- []byte) bool {
	ticker := time.NewTicker(time.Second)
	scopelog.Printf(scopeConsumer, "init partitionConsumer\n")
	partitionConsumer, err := consumer.ConsumePartition(c.Topic, int32(c.Partition), c.getOffset())
	if err != nil {
		scopelog.Printf(scopeConsumer, "Get ConsumePartition err:%v\n", err)
		return false
	}
	defer partitionConsumer.Close()

	var lastReceivedMsgTime time.Time
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			if c.isReplay {
				c.isReplay = false
			}
			c.Offset = &msg.Offset

			lastReceivedMsgTime = time.Now()
			scopelog.Printf(scopeConsumer, "offset: [%d], message: [%s],  lastReceivedMsgTime:%v\n", msg.Offset, msg.Value, lastReceivedMsgTime)
			scopelog.Printf(scopeConsumer, "predict, ip:%s, host: %s, auth:%s ,timeout: %d\n", c.Options.PredictIp, c.Options.PredictDomain, c.Options.PredictAuth, c.Options.PredictTimeout)
			//TODO remove, just for test
			data := []byte(`{"instances": [{"x1":6.2, "x2":2.2, "x3":1.1, "x4":1.2}]}`)

			result, err := predict.Post(c.Options.PredictIp,
				c.Options.PredictDomain,
				c.Options.PredictAuth,
				c.Options.PredictTimeout,
				data)
			if err != nil {
				result = "predict error"
			}
			zipData := fmt.Sprintf("origin data: %s | prediction: %s", msg.Value, result)
			//_ = wsConn.WriteMessage(websocket.TextMessage, msg.Value)
			dataCh <- []byte(zipData)

		case <-ticker.C:
			fmt.Printf("check if replay from oldest, lastTime:%v, 时间差: %v\n", lastReceivedMsgTime, time.Now().Sub(lastReceivedMsgTime))
			if time.Now().After(lastReceivedMsgTime.Add(time.Minute)) {
				c.isReplay = true
				return true
			}
		case <-c.StopCh:
			scopelog.Printf(scopeConsumer, "Consumer Stopped (Precisely)\n")
			return false
		}
		time.Sleep(2 * time.Second)
	}

}
