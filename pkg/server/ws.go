package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"opcdata-predict-client/cmd/option"
	"opcdata-predict-client/pkg/kafka"
	"opcdata-predict-client/pkg/scopelog"
	"sync"
)

type WebSocketManager struct {
	lock         sync.RWMutex
	consumeState bool
	WsConn       *websocket.Conn
	Options      option.Options
}

func NewWebSocketManager(opts option.Options) *WebSocketManager {
	m := &WebSocketManager{}
	m.consumeState = false
	m.lock = sync.RWMutex{}
	m.Options = opts
	return m
}

var myUpGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var scopeWsHandler = "WsHandler"

func (m *WebSocketManager) WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := myUpGrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		//log.Println("升级为websocket失败", err.Error())
		scopelog.Printf(scopeWsHandler, "升级为websocket失败,err: %v\n", err)
		return
	}
	m.WsConn = conn
	defer conn.Close()

	scopelog.Printf(scopeWsHandler, "Init Websocket Connection\n")
	consumeMgr := kafka.NewMessageCollector()
	scopelog.Printf(scopeWsHandler, "Init Kafka Message Collector\n")
	m.CommandHandler(consumeMgr)
}

var scopeCommandHandler = "CommandHandler"
var actionRequest = "User Request"
var actionState = "State"

func (m *WebSocketManager) CommandHandler(consumerMgr *kafka.MessageCollector) {
	log.Printf("[%s]: trap in", scopeCommandHandler)
	for {

		// Read message from browser
		msgType, msg, err := m.WsConn.ReadMessage()
		if err != nil {
			scopelog.Printf(scopeCommandHandler, "ReadMessage err: %v\n", err)
			return
		}
		scopelog.Printf(scopeCommandHandler, "msgType: %v\n", msgType)
		if msgType != websocket.TextMessage {
			continue
		}
		scopelog.Printf(scopeCommandHandler, "Receive command\n")
		scopelog.Printf(scopeCommandHandler, "msg: %v\n", msg)
		if string(msg) == ControlCommand.Start {
			var state bool
			m.lock.RLock()
			state = m.consumeState
			m.lock.RUnlock()

			if state == false {
				scopelog.Printf(actionRequest, "StartUp Kafka Consumer\n")
				if err := consumerMgr.CreateClient(); err == nil {
					if err := consumerMgr.StartConsume(m.WsConn,m.Options); err != nil {
						scopelog.Printf(actionState, "Consumer Start err, %v\n", err)
					} else {
						scopelog.Printf(actionState, "Consumer Started\n")
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
				scopelog.Printf(actionRequest, "Stop Kafka Consumer")
				consumerMgr.StopConsume()
				scopelog.Printf(actionState, "Consumer stopped\n")
				// set state
				m.lock.Lock()
				m.consumeState = false
				m.lock.Unlock()
			}

		}

	}
}
