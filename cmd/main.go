package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"net/http"
	"opcdata-predict-client/pkg/server"
)

var (
	// predict
	host     string
	auth     string
	serverIp string
	// kafka
	//command           string
	brokers           string
	topic             string
	partition         int
	saslEnable        bool
	username          string
	password          string
	saslSASLMechanism string
)

func main() {
	setupCli()
	fmt.Println(host)
	var predictResult = make(chan []byte)

	go func() {
		i := 0
		for {
			if i > 10 {
				break
			}
			predictResult <- []byte(fmt.Sprintf("hello - %d", i))

			i = i + 1
		}
	}()

	webSocketServerStart(predictResult)

	/*
		if host == "" || auth == "" || serverIp == "" {
			fmt.Println("--host or --auth or --server is missing")
			os.Exit(1)
		}

		// websocket server
		wssSvr := server.WsServer{}
		svr := wssSvr.Create(":8080")
		go func() {
			fmt.Println("listen to :8080")
			if err := svr.ListenAndServe(); err != nil {
				panic(err)
			}
		}()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop

		wssSvr.Shutdown(30)
	*/

}

//var lock sync.Mutex
//var consumeState bool
func webSocketServerStart(resultChan <-chan []byte) {
	//var upGrader = websocket.Upgrader{
	//	ReadBufferSize:  1024,
	//	WriteBufferSize: 1024,
	//}
	/*
	   	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
	   		conn, err := upGrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	   		if err != nil {
	   			log.Println("升级为websocket失败", err.Error())
	   			return
	   		}
	   		defer conn.Close()
	   		wsmgr := server.NewWebSocketManager()
	   		go wsmgr.CommandHandler(conn)
	   		//time.Sleep(10*time.Second)
	   		/*
	   			ticker := time.NewTicker(400 * time.Millisecond)
	   			for {
	   				select {
	   				case byteArr := <-resultChan:
	   					conn.WriteMessage(websocket.TextMessage, byteArr)

	   				case <-ticker.C:
	   					// nothing to do
	   				}
	   					//// Read message from browser
	   					//msgType, msg, err := conn.ReadMessage()
	   					//if err != nil {
	   					//	return
	   					//}
	   					//// Print the message to the console
	   					//fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
	   					//fmt.Printf("msgType=%v", msgType)
	   					//// Write message back to browser
	   					//if err := conn.WriteMessage(msgType, msg); err != nil {
	   					//	return
	   					//}

	   			}
	   /
	   	})
	*/
	wsmgr := server.NewWebSocketManager()
	http.HandleFunc("/echo", wsmgr.WsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	fmt.Println("Listening :8080")
	http.ListenAndServe(":8080", nil)
}

func setupCli() {
	flag.StringVar(&host, "host", "", "model domain name")
	flag.StringVar(&auth, "auth", "", "model auth code")
	flag.StringVar(&serverIp, "server", "", "server ip that will be request")

	flag.StringVar(&brokers, "brokers", "localhost:9092", "Common separated kafka hosts")
	flag.StringVar(&topic, "topic", "test", "Kafka topic")
	flag.IntVar(&partition, "partition", 0, "Kafka topic partition")

	flag.BoolVar(&saslEnable, "sasl", false, "SASL enable")
	flag.StringVar(&username, "username", "", "SASL Username")
	flag.StringVar(&password, "password", "", "SASL Password")
	flag.StringVar(&saslSASLMechanism, "mechanism", sarama.SASLTypePlaintext, "SASL mechanism")

	flag.Parse()
}
