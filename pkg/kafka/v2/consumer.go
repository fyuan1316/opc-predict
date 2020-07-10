package v2

import (
	"github.com/Shopify/sarama"
	"opcdata-predict/cmd/option"
	"opcdata-predict/pkg/scopelog"
	"opcdata-predict/pkg/server"
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

func NewKafkaConsumer(opts option.Options) *KafkaConsumer {
	c := &KafkaConsumer{}
	c.Brokers = opts.Brokers     //"localhost:9092"
	c.Topic = opts.Topic         //"test"
	c.Partition = opts.Partition // 0
	c.Config = sarama.NewConfig()
	if opts.SaslEnable {
		c.Config.Net.SASL.Enable = true
		c.Config.Net.SASL.User = opts.SaslUsername
		c.Config.Net.SASL.Password = opts.SaslPassword
		c.Config.Net.SASL.Mechanism = sarama.SASLMechanism(opts.SaslMechanism)
	}
	c.StopCh = make(chan struct{})
	return c
}

var (
	doOnce        sync.Once
	scopeConsumer = "Consumer"
)

//TODO use sarama.OffsetOldest
var beginOffset = sarama.OffsetOldest //565

func (c *KafkaConsumer) CreateClient() error {
	var client sarama.Client
	var err error
	doOnce.Do(func() {
		client, err = sarama.NewClient(strings.Split(c.Brokers, ","), c.Config)
		p := int64(beginOffset)
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
	if err != nil {
		return err
	}
	//p := sarama.OffsetOldest
	p := int64(beginOffset)
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
			scopelog.Printf(scopeConsumer, "replay")
		}
	}()
	return nil
}

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
		dataCh <- []byte(server.CommandPrefix + err.Error())
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
			dataCh <- msg.Value
		case <-ticker.C:
			scopelog.Printf(scopeConsumer, "check if replay from oldest, lastTime:%v, 时间差: %v\n", lastReceivedMsgTime, time.Now().Sub(lastReceivedMsgTime))
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
