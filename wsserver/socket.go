package wsserver

import (
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

func (s *socket) receive(msgHandler func(int, []byte), errorHandler func(error)) {
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
			return
		case msg := <-msgChan:
			if msgHandler != nil {
				msgHandler(s.id, msg)
			}
		}
	}
}

func (s *socket) send(errHandler func(error)) {
	for {
		select {
		case <-s.stopSignal:
			return

		case msg := <-s.sendMsgChan:
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

func (server *WsServer) querySockets(ids ...int) []*socket {
	var socks map[int]*socket
	server.sockets.BatchQuery(ids, &socks)
	var sockets []*socket
	for _, v := range socks {
		sockets = append(sockets, v)
	}
	return sockets
}

func (server *WsServer) ListSockets() []int {
	var ids []int
	smap := server.sockets.Interface().(socketsMap)
	for k := range smap {
		ids = append(ids, k)
	}
	return ids
}
