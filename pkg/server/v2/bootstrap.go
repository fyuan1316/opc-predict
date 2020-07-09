package v2

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"opcdata-predict/cmd/option"
	kafka2 "opcdata-predict/pkg/kafka/v2"
	"opcdata-predict/pkg/predict"
)

type WsServer struct {
	Port    int
	Options option.Options
	//
	KafkaManager   *kafka2.ConsumerManager
	PredictService *predict.Manager
}

func NewWsServer(port int, opts option.Options) WsServer {
	s := WsServer{}
	s.Port = port
	s.Options = opts
	s.KafkaManager = kafka2.NewConsumerManager()
	s.PredictService = predict.NewManager(s.Options)

	return s
}

func (s *WsServer) Boot() {
	//
	go clientsManager.StartMessageLoop()                // multiplex sockets
	go s.KafkaManager.Process(clientsManager.Broadcast) // kafka worker

	//route
	http.HandleFunc("/echo", s.myWsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	log.Printf("Listening :%d\n", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}

var myUpGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WsServer) myWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := myUpGrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	//注册当前的client
	client := &Client{Conn: conn, Send: make(chan []byte)}
	clientsManager.Register <- client

	// 接收client 发来的命令，写入 kafka mgr的命令管道
	go client.Read(s.KafkaManager.CommandCh)
	// 将数据写会客户端
	go client.Write(s.PredictService.GetPredictedResult)

}
