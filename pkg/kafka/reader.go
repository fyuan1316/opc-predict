package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"strings"
	"sync"
	"time"
)

type MessageCollector struct {
	Config    *sarama.Config
	Brokers   string
	Topic     string
	Partition int
	Client    sarama.Client
	Offset    *int64
	StopCh    chan struct{}
}

func NewMessageCollector() *MessageCollector {
	c := &MessageCollector{}
	c.Brokers = "localhost:9092"
	c.Topic = "test"
	c.Partition = 0
	c.Config = sarama.NewConfig()
	c.StopCh = make(chan struct{})
	return c
}

var doOnce sync.Once

func (c *MessageCollector) CreateClient() error {
	var client sarama.Client
	var err error
	doOnce.Do(func() {
		client, err = sarama.NewClient(strings.Split(c.Brokers, ","), c.Config)
		p := sarama.OffsetOldest
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

func (c *MessageCollector) StopConsume() {
	close(c.StopCh)
}

func (c *MessageCollector) StartConsume() error {
	consumer, err := sarama.NewConsumerFromClient(c.Client)
	if err != nil {
		return err
	}
	//if c.StopCh == nil {
		c.StopCh = make(chan struct{})
	//}
	go func() {
		for {
			if replay := c.consume(consumer); !replay {
				break
			}
			fmt.Println("replay")
		}
	}()
	return nil
}

func (c *MessageCollector) consume(consumer sarama.Consumer) bool {
	ticker := time.NewTicker(time.Second)
	fmt.Println("init partitionConsumer")
	partitionConsumer, err := consumer.ConsumePartition(c.Topic, int32(c.Partition), *c.Offset)
	if err != nil {
		log.Println(err)
		return false
	}
	defer partitionConsumer.Close()

	//for {
	//	msg := <-partitionConsumer.Messages()
	//	log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)
	//}
	var lastReceivedMsgTime time.Time
	for {


		select {
		case msg := <-partitionConsumer.Messages():
			lastReceivedMsgTime = time.Now()
			log.Printf("Consumed message: [%s], offset: [%d], lastReceivedMsgTime:%v\n", msg.Value, msg.Offset, lastReceivedMsgTime)
			fmt.Println("todo call prediction")

			//
			c.Offset = &msg.Offset
		case <-ticker.C:
			fmt.Printf("check if replay from oldest, lastTime:%v, 时间差: %v\n", lastReceivedMsgTime, time.Now().Sub(lastReceivedMsgTime))
			if time.Now().After(lastReceivedMsgTime.Add(time.Minute)) {
				//goto ConsumeFromOldestOffset
				return true
			}
		case <-c.StopCh:
			fmt.Println("user stop consume")
			goto breakout
		}
		time.Sleep(time.Second)
	}
breakout:
	return false
}
