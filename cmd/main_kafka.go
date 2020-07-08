package cmd


/*
import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	command           string
	hosts             string
	topic             string
	partition         int
	saslEnable        bool
	username          string
	password          string
	saslSASLMechanism string
	tlsEnable         bool
	clientcert        string
	clientkey         string
	cacert            string
)

func main() {
	flag.StringVar(&command, "command", "consumer", "consumer|producer")
	flag.StringVar(&hosts, "host", "localhost:9092", "Common separated kafka hosts")
	flag.StringVar(&topic, "topic", "test", "Kafka topic")
	flag.IntVar(&partition, "partition", 0, "Kafka topic partition")

	flag.BoolVar(&saslEnable, "sasl", false, "SASL enable")
	flag.StringVar(&username, "username", "", "SASL Username")
	flag.StringVar(&password, "password", "", "SASL Password")
	flag.StringVar(&saslSASLMechanism, "mechanism", sarama.SASLTypePlaintext, "SASL mechanism")

	flag.BoolVar(&tlsEnable, "tls", false, "TLS enable")
	flag.StringVar(&clientcert, "cert", "cert.pem", "Client Certificate")
	flag.StringVar(&clientkey, "key", "key.pem", "Client Key")
	flag.StringVar(&cacert, "ca", "ca.pem", "CA Certificate")
	flag.Parse()

	config := sarama.NewConfig()
	if saslEnable {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.Mechanism = sarama.SASLMechanism(saslSASLMechanism)
	}
	fmt.Println("Mechanism=", config.Net.SASL.Mechanism)

	if tlsEnable {
		//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
		tlsConfig, err := genTLSConfig(clientcert, clientkey, cacert)
		if err != nil {
			log.Fatal(err)
		}

		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	client, err := sarama.NewClient(strings.Split(hosts, ","), config)
	if err != nil {
		log.Fatalf("unable to create kafka client: %q", err)
	}

	if command == "consumer" {
		consumer, err := sarama.NewConsumerFromClient(client)
		if err != nil {
			log.Fatal(err)
		}
		defer consumer.Close()
		loopConsumer(consumer, topic, partition)
	} else {
		producer, err := sarama.NewAsyncProducerFromClient(client)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()
		loopProducer(producer, topic, partition)
	}

}
func loopProducer(producer sarama.AsyncProducer, topic string, partition int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
		} else if text == "exit" || text == "quit" {
			break
		} else {
			producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(text)}
			log.Printf("Produced message: [%s]\n", text)
		}
		fmt.Print("> ")
	}
}

func loopConsumer(consumer sarama.Consumer, topic string, partition int) {

	for {
		if replay := loopConsumerWorker(consumer, topic, partition); !replay {
			return
		}
		fmt.Println("replay")
	}
}

func loopConsumerWorker(consumer sarama.Consumer, topic string, partition int) bool {
	//partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	//partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), 565)

	ticker := time.NewTicker(time.Second)
	//ConsumeFromOldestOffset:
	fmt.Println("init partitionConsumer")
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), 565)
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
			//fmt.Println(msg, lastReceivedMsgTime)
			fmt.Println("todo call prediction")
		case <-ticker.C:
			fmt.Printf("check if replay from oldest, lastTime:%v, 时间差: %v\n", lastReceivedMsgTime, time.Now().Sub(lastReceivedMsgTime))
			if time.Now().After(lastReceivedMsgTime.Add(time.Minute)) {
				//goto ConsumeFromOldestOffset
				return true
			}
		}
		time.Sleep(time.Second)
	}

}
func genTLSConfig(clientcertfile, clientkeyfile, cacertfile string) (*tls.Config, error) {
	// load client cert
	clientcert, err := tls.LoadX509KeyPair(clientcertfile, clientkeyfile)
	if err != nil {
		return nil, err
	}

	// load ca cert pool
	cacert, err := ioutil.ReadFile(cacertfile)
	if err != nil {
		return nil, err
	}
	cacertpool := x509.NewCertPool()
	cacertpool.AppendCertsFromPEM(cacert)

	// generate tlcconfig
	tlsConfig := tls.Config{}
	tlsConfig.RootCAs = cacertpool
	tlsConfig.Certificates = []tls.Certificate{clientcert}
	tlsConfig.BuildNameToCertificate()
	// tlsConfig.InsecureSkipVerify = true // This can be used on test server if domain does not match cert:
	return &tlsConfig, err
}
*/