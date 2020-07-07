package server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type WsServer struct {
	Addr   string
	server http.Server
}

func (s *WsServer) Create(addr string) http.Server {
	s.Addr = addr
	s.server = http.Server{}
	addHandlers()
	return s.server
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func addHandlers() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upGrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			//msg = []byte("echo"+string(msg))
			msg = []byte("fangyuan")
			if err := conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
}

type wssHandler struct {
}

func (h wssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (s *WsServer) Shutdown(waitSeconds int) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(waitSeconds)*time.Second)
	if err := s.server.Shutdown(ctx); err != nil {
		// handle err
		log.Printf("shutdown websocket server err: %v\n", err)
	}
}
