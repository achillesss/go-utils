package wsserver

import (
	"github.com/achillesss/go-utils/log"
	"golang.org/x/net/websocket"
)

type socketsMap map[int]*socket

type socket struct {
	id           int
	conn         *websocket.Conn
	sendMsgChan  chan []byte
	stopSignal   chan struct{}
	isConnecting *bool
}

func (s *socket) stop() {
	if *s.isConnecting {
		*s.isConnecting = false
		close(s.stopSignal)
	}
}

func (s *socket) receive(msgHandler func([]byte), errorHandler func(error)) {
	log.Infofln("start receive")
	msgChan := make(chan []byte)
	go func() {
		for {
			var msg []byte
			err := websocket.Message.Receive(s.conn, &msg)
			switch err {
			case nil:
				msgChan <- msg
			default:
				if errorHandler != nil {
					errorHandler(err)
				}
				s.stop()
				return
			}
		}
	}()

	for {
		select {
		case <-s.stopSignal:
			log.Infofln("receive stop signal")
			return
		case msg := <-msgChan:
			if msgHandler != nil {
				msgHandler(msg)
			}
		}
	}
}

func (s *socket) send(errHandler func(error)) {
	log.Infofln("start send")
	for {
		select {
		case <-s.stopSignal:
			log.Infofln("receive stop signal")
			return

		case msg := <-s.sendMsgChan:
			log.Infofln("send %s", msg)
			err := websocket.Message.Send(s.conn, msg)
			if err != nil {
				if errHandler != nil {
					errHandler(err)
				}
				s.stop()
			}
		}
	}
}

func (server *WsServer) newSocket(conn *websocket.Conn) *socket {
	var s socket
	*server.lastConnID++
	s.id = *server.lastConnID
	s.isConnecting = new(bool)
	*s.isConnecting = true
	s.conn = conn
	s.sendMsgChan = make(chan []byte)
	s.stopSignal = make(chan struct{})
	return &s
}

func (server WsServer) querySocket(id int) *socket {
	var s *socket
	server.sockets.Query(id, &s)
	return s
}

func (server WsServer) ListSockets() []int {
	var ids []int
	smap := server.sockets.Interface().(socketsMap)
	for k := range smap {
		ids = append(ids, k)
	}
	return ids
}
