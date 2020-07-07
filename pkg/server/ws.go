package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"opcdata-predict-client/pkg/kafka"
	"sync"
)

type WebSocketManager struct {
	lock         sync.RWMutex
	consumeState bool
	WsConn       *websocket.Conn
}

func NewWebSocketManager() *WebSocketManager {
	m := &WebSocketManager{}
	m.consumeState = false
	m.lock = sync.RWMutex{}
	return m
}

var myUpGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (m *WebSocketManager) WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := myUpGrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		log.Println("升级为websocket失败", err.Error())
		return
	}
	m.WsConn = conn
	defer conn.Close()

	fmt.Println("hahahah")
	consumeMgr := kafka.NewMessageCollector()

	m.CommandHandler(consumeMgr)
}

func (m *WebSocketManager) CommandHandler(consumerMgr *kafka.MessageCollector) {
	fmt.Println("222323 hahah ")
	for {
		fmt.Println("fangfangfang ")
		// Read message from browser
		msgType, msg, err := m.WsConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("msgType: %v\n", msgType)
		if msgType != websocket.TextMessage {
			continue
		}
		fmt.Printf("msg: %v\n", msg)
		if string(msg) == ControlCommand.Start {
			var state bool
			m.lock.RLock()
			state = m.consumeState
			m.lock.RUnlock()

			if state == false {
				// startup consumer TODO
				fmt.Println("startup consumer")
				if err := consumerMgr.CreateClient(); err == nil {
					if err := consumerMgr.StartConsume(); err != nil {
						log.Printf("StartConsume err: %v\n", err)
					}
				}
				// set state
				m.lock.Lock()
				m.consumeState = true
				m.lock.Unlock()
			}

		}
		if string(msg) == ControlCommand.Stop {
			var state bool
			m.lock.RLock()
			state = m.consumeState
			m.lock.RUnlock()

			if state == true {
				// stop consumer TODO
				fmt.Println("stop consumer")
				consumerMgr.StopConsume()
				fmt.Println("stopped consumer 2")
				// set state
				m.lock.Lock()
				m.consumeState = false
				m.lock.Unlock()
			}

		}

	}
}
