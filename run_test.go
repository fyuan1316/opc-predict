package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	kafka2 "opcdata-predict/pkg/kafka/v2"
	"opcdata-predict/pkg/predict"
	server2 "opcdata-predict/pkg/server/v2"
	"testing"
)

func Test_Run(t *testing.T) {

	fmt.Println("listen :8080")
	Start()
}

var kafkaMgr = kafka2.NewConsumerManager()
var predictService  = predict.NewManager()
func Start() {
	fmt.Println("start")
	go server2.ClientManager.StartMessageLoop()
	go kafkaMgr.Process(server2.ClientManager.Broadcast)
	http.HandleFunc("/echo", myWsHandler)
	http.ListenAndServe(":8080", nil)
}

var myUpGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func myWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := myUpGrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	//注册当前的client
	client := &server2.Client{Conn: conn, Send: make(chan []byte)}
	server2.ClientManager.Register <- client

	// 接收client 发来的命令，写入 kafka mgr的命令管道
	go client.Read(kafkaMgr.CommandCh)
	// 将数据写会客户端
	go client.Write()

}
