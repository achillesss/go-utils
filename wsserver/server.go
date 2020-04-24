package wsserver

import (
	"net"
	"net/http"

	"github.com/achillesss/go-utils/log"
	"golang.org/x/net/websocket"
)

func (server *WsServer) onConnection(c *websocket.Conn) {
	s := server.newSocket(c)

	go func() {
		<-s.stopSignal

		server.socketMsgChanLock.Lock()
		delete(server.socketMsgChans, s.id)
		server.socketMsgChanLock.Unlock()

		close(s.sendMsgChan)

		server.socketStopChanLock.Lock()
		delete(server.socketStopChans, s.id)
		server.socketStopChanLock.Unlock()

	}()

	go s.send(server.sendMsgHandler, server.sendErrorHandler)
	s.receive(
		func(msg []byte) {
			if server.ShouldAnswerPing && string(msg) == "ping" {
				s.sendMsgChan <- []byte("pong")
			}
			server.receiveMsgHandler(msg)
		},
		server.receiveErrorHandler,
	)
}

func NewWsServer(addr, pattern string, options ...ServerOption) *WsServer {
	var server WsServer
	server.routers = make(map[string]map[string]http.Handler)
	server.socketMsgChans = make(map[int64]chan []byte)
	server.socketStopChans = make(map[int64]chan struct{})
	options = append(options, WithRouter(addr, pattern, ConvertWsHandlerFunctionToHttpHandler(server.onConnection)))
	for _, option := range options {
		option.updateOption(&server)
	}

	// 	server.rooms = gomap.NewMap(make(roomsMap))
	// 	server.roomMsg = make(chan map[int64][]byte, 1)

	// server.initRoomsMap()
	// server.initChannels()
	// go server.chatMonitor()
	return &server
}

func WithSendErrorHandler(f func(error)) ServerOption {
	return newOptionHolder(func(ws *WsServer) {
		ws.sendErrorHandler = f
	})
}

func WithRecvErrorHandler(f func(error)) ServerOption {
	return newOptionHolder(func(ws *WsServer) {
		ws.receiveErrorHandler = f
	})
}

func WithSendMsgHandler(f func([]byte)) ServerOption {
	return newOptionHolder(func(ws *WsServer) {
		ws.sendMsgHandler = f
	})
}

func WithRecvMsgHandler(f func([]byte)) ServerOption {
	return newOptionHolder(func(ws *WsServer) {
		ws.receiveMsgHandler = f
	})
}

func WithRouter(addr, pattern string, handler http.Handler) ServerOption {
	return newOptionHolder(func(ws *WsServer) {
		var r = ws.routers[addr]
		if r == nil {
			r = make(map[string]http.Handler)
			ws.routers[addr] = r
		}
		r[pattern] = handler
	})
}

func WithAnswerPing() ServerOption {
	return newOptionHolder(func(ws *WsServer) {
		ws.ShouldAnswerPing = true
	})
}

func (server *WsServer) Serve() {
	for addr, r := range server.routers {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			panic(err)
		}

		var mux = http.NewServeMux()
		for pattern, handler := range r {
			log.Infofln("new router: %s", pattern)
			mux.Handle(pattern, handler)
		}

		go http.Serve(l, mux)
		log.Infofln("serve above routers on %s", addr)
	}
}

func (server *WsServer) SendMsgTo(msg []byte, ids ...int64) {
	server.socketMsgChanLock.RLock()
	defer server.socketMsgChanLock.RUnlock()

	for _, id := range ids {
		var ch, ok = server.socketMsgChans[id]
		if !ok {
			continue
		}

		go func(ch chan []byte) {
			select {
			case <-ch:
			default:
				ch <- msg
			}
		}(ch)
	}
}

func (server *WsServer) Broadcast(msg []byte) {
	server.socketMsgChanLock.RLock()
	defer server.socketMsgChanLock.RUnlock()

	for _, ch := range server.socketMsgChans {
		go func(ch chan []byte) {
			select {
			case <-ch:
			default:
				ch <- msg
			}
		}(ch)
	}
}

// func (server *WsServer) SendRoomMsg(msg []byte, roomID int64) {
// 	server.roomMsg <- map[int64][]byte{roomID: msg}
// }

func (server *WsServer) Stop(socketID int64) {
	server.socketStopChanLock.RLock()
	defer server.socketStopChanLock.RUnlock()

	var ch, ok = server.socketStopChans[socketID]
	if !ok {
		return
	}

	select {
	case <-ch:
	default:
		close(ch)
	}
}
