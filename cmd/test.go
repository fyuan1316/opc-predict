package cmd

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var (
	testhosts             string
	testtopic             string
	testpartition         int
	testsaslEnable        bool
	testusername          string
	testpassword          string
	testsaslSASLMechanism string
	//tlsEnable         bool
	//clientcert        string
	//clientkey         string
	//cacert            string
)
var consumeCommand = &cobra.Command{
	Use:   "consume-sasl",
	Short: "read data from topic with sasl",
	Long:  `read data from topic with sasl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consume called")
		RunCommand()
	},
}

func init() {
	initflags()
	rootCmd.AddCommand(consumeCommand)

}

func initflags() {

	//kafka
	consumeCommand.Flags().StringVar(&testhosts, "brokers", "localhost:9092", "kafka's brokers list by comma split")
	consumeCommand.Flags().StringVar(&testtopic, "topic", "test", "kafka topic name to consume")
	consumeCommand.Flags().IntVar(&testpartition, "partition", 0, "kafka partition")
	consumeCommand.Flags().BoolVar(&testsaslEnable, "sasl", false, "sasl switcher")
	consumeCommand.Flags().StringVar(&testusername, "username", "", "sasl username")
	consumeCommand.Flags().StringVar(&testpassword, "passwd", "", "sasl password")
	consumeCommand.Flags().StringVar(&testsaslSASLMechanism, "mechanism", "PLAIN", "sasl mechanism")

}

func RunCommand() {
	config := sarama.NewConfig()
	if testsaslEnable {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = testusername
		config.Net.SASL.Password = testpassword
		config.Net.SASL.Mechanism = sarama.SASLMechanism(testsaslSASLMechanism)
	}
	fmt.Println("Mechanism=", config.Net.SASL.Mechanism)
	fmt.Println("brokers=", testhosts)
	fmt.Println("topic=", testtopic)
	fmt.Println("partition=", testpartition)
	fmt.Println("testsaslEnable=", testsaslEnable)
	fmt.Println("testusername=", testusername)
	fmt.Println("testpassword=", testpassword)
	fmt.Println("testsaslSASLMechanism=", testsaslSASLMechanism)

	//if tlsEnable {
	//	//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	//	tlsConfig, err := genTLSConfig(clientcert, clientkey, cacert)
	//	tlsConfig, err := genTLSConfig(clientcert, clientkey, cacert)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	config.Net.TLS.Enable = true
	//	config.Net.TLS.Config = tlsConfig
	//}

	client, err := sarama.NewClient(strings.Split(testhosts, ","), config)
	if err != nil {
		log.Fatalf("unable to create kafka client: %q", err)
	}

	//if command == "consumer" {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	loopConsumer(consumer, testtopic, testpartition)
	//} else {
	//	producer, err := sarama.NewAsyncProducerFromClient(client)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer producer.Close()
	//	loopProducer(producer, topic, partition)
	//}
}
func loopConsumer(consumer sarama.Consumer, topic string, partition int) {
	//partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	//partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), 565)

	fmt.Println("init partitionConsumer")
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer partitionConsumer.Close()

	for {
		msg := <-partitionConsumer.Messages()
		time.Sleep(time.Second)
		log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)
	}
}
