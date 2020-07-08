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
//
//func (s *WsServer) Shutdown(waitSeconds int) {
//	ctx, _ := context.WithTimeout(context.Background(), time.Duration(waitSeconds)*time.Second)
//	if err := s.server.Shutdown(ctx); err != nil {
//		// handle err
//		log.Printf("shutdown websocket server err: %v\n", err)
//	}
//}
