package v2

import (
	"fmt"
	"github.com/gorilla/websocket"
	"opcdata-predict/pkg/server"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}


/**
将从client 读取到的消息 发送到 命令通道中，这里要识别出是否合法的命令
*/
func (c *Client) Read(commandCh chan<- []byte) {
	defer func() {
		clientsManager.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			clientsManager.Unregister <- c
			c.Conn.Close()
			break
		}
		fmt.Printf("client recv: %v\n", message)
		if isValideCommand(message) {
			commandCh <- message
		} else {
			fmt.Printf("invalidate command[%s], skipped\n", message)
		}

		//jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		//manager.broadcast <- jsonMessage
	}
}

func (c *Client) Write(predictService func([]byte)[]byte) {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.WriteMessage(websocket.TextMessage, predictService(msg))
		}
	}
}

func isValideCommand(c []byte) bool {
	switch string(c) {
	case server.ControlCommand.Stop, server.ControlCommand.Start:
		return true
	}
	return false
}

type ClientsManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	//lock       sync.RWMutex
}

var clientsManager = ClientsManager{
	Clients:    make(map[*Client]bool),
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (m *ClientsManager) StartMessageLoop() {
	for {
		select {
		case conn := <-m.Register:
			//m.lock.Lock()
			m.Clients[conn] = true
			//m.lock.Unlock()
		case conn := <-m.Unregister:
			//m.lock.Lock()
			if _, ok := m.Clients[conn]; ok {
				delete(m.Clients, conn)
			}
			//m.lock.Unlock()
		case msg := <-m.Broadcast:
			fmt.Printf("ClientsManager recv: %s\n", msg)
			for conn := range m.Clients {
				//connection 正常时，msg 发送回client 自己的消息chan；异常时将client 清理掉.
				select {
				case conn.Send <- msg:
				default:
					close(conn.Send)
					delete(m.Clients, conn)
				}
			}
		}

	}
}
