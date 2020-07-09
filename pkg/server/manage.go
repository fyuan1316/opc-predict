package server

import (
	"fmt"
	"log"
	"net/http"
	"opcdata-predict/cmd/option"
)

type WsServer struct {
	Port    int
	Options option.Options
}

func NewWsServer(port int, opts option.Options) WsServer {
	s := WsServer{}
	s.Port = port
	s.Options = opts
	return s
}

func (s *WsServer) Boot() {
	wsMgr := NewWebSocketManager(s.Options)
	http.HandleFunc("/echo", wsMgr.WsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	log.Printf("Listening :%d\n", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}
