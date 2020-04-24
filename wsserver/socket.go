package wsserver

import (
	"sync/atomic"

	"golang.org/x/net/websocket"
)

type socketsMap map[int64]*socket

type socket struct {
	id          int64
	conn        *websocket.Conn
	sendMsgChan chan []byte
	stopSignal  chan struct{}
}

func (s *socket) try(f func()) {
	select {
	case <-s.stopSignal:
	default:
		f()
	}
}

func (s *socket) tryLoop(f func()) {
	for {
		select {
		case <-s.stopSignal:
			return

		default:
			f()
		}
	}
}

func (s *socket) stop() {
	s.try(func() {
		close(s.stopSignal)
	})
}

func (s *socket) receive(msgHandler MsgHandler, errHandler ErrHandler) {
	go func() {
		for {
			var msg []byte
			var err = websocket.Message.Receive(s.conn, &msg)
			if err != nil {
				if errHandler != nil {
					go errHandler(err)
				}
				s.stop()
				return
			}

			if msgHandler != nil {
				go msgHandler(msg)
			}
		}
	}()

	select {
	case <-s.stopSignal:
	}
}

func (s *socket) sendMsg(msg []byte, msgHandler MsgHandler, errHandler ErrHandler) {
	if msgHandler != nil {
		go msgHandler(msg)
	}

	var err = websocket.Message.Send(s.conn, msg)
	if err == nil {
		return
	}

	if errHandler != nil {
		go errHandler(err)
	}

	s.stop()
}

func (s *socket) send(msgHandler MsgHandler, errHandler ErrHandler) {
	for {
		select {
		case <-s.stopSignal:
			return

		case msg := <-s.sendMsgChan:
			go s.sendMsg(msg, msgHandler, errHandler)
		}
	}
}

func (server *WsServer) newSocket(conn *websocket.Conn) *socket {
	var s socket
	s.id = atomic.AddInt64(&server.lastConnID, 1)
	s.conn = conn
	s.sendMsgChan = make(chan []byte)
	s.stopSignal = make(chan struct{})

	server.socketMsgChanLock.Lock()
	defer server.socketMsgChanLock.Unlock()
	server.socketMsgChans[s.id] = s.sendMsgChan

	server.socketStopChanLock.Lock()
	defer server.socketStopChanLock.Unlock()
	server.socketStopChans[s.id] = s.stopSignal

	return &s
}

func (server *WsServer) ListSocketIDs() []int64 {
	server.socketMsgChanLock.RLock()
	defer server.socketMsgChanLock.RUnlock()

	var ids []int64
	for k := range server.socketMsgChans {
		ids = append(ids, k)
	}
	return ids
}
