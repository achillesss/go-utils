package wsserver

import (
	"fmt"
	"net"

	"github.com/achillesss/go-utils/go-map"
	"github.com/achillesss/go-utils/log"
	"golang.org/x/net/websocket"
)

func (server *WsServer) initSocketsMap() {
	server.sockets = gomap.NewMap(make(socketsMap))
}

func (server *WsServer) initRoomsMap() {
	server.rooms = gomap.NewMap(make(roomsMap))
}

func (server *WsServer) initChannels() {
	server.roomMsg = make(chan map[int][]byte, 1)
}

func (server *WsServer) onConnection(c *websocket.Conn) {
	log.Infofln("new connection")
	s := server.newSocket(c)
	server.sockets.Add(s.id, s)
	go s.send(server.sendErrorHandler)
	s.receive(
		func(socketID int, msg []byte) {
			if server.ShouldAnswerPing && string(msg) == "ping" {
				s.sendMsgChan <- []byte("pong")
			}
			server.receiveMsgHandler(socketID, msg)
		},
		server.receiveErrorHandler,
	)
}

func NewWsServerFromConfig(config *wsServerConfig) *WsServer {
	var server WsServer
	server.lastConnID = new(int)
	server.lastRoomID = new(int)
	server.handlers = config.handlers
	server.initSocketsMap()
	server.initRoomsMap()
	server.initChannels()
	go server.chatMonitor()
	return &server
}

func (server *WsServer) SetSendErrorHandler(f func(error)) {
	server.sendErrorHandler = f
}

func (server *WsServer) SetReceiveErrorHandler(f func(error)) {
	server.receiveErrorHandler = f
}

func (server *WsServer) SetReceiveMsgHandler(f func(int, []byte)) {
	server.receiveMsgHandler = f
}

func (server WsServer) SetConnectionRouters(addr net.Addr, patterns ...string) {
	h := newWsHandler()
	h.address = addr
	for _, pattern := range patterns {
		h.registerRouter(pattern, ConvertWsHandlerFunctionToHttpHandler(server.onConnection))
	}
	h.serve()
}

func (server WsServer) Serve() {
	for _, h := range server.handlers {
		h.serve()
	}
}

func (server WsServer) SendTo(msg []byte, socketID int) error {
	s := server.querySocket(socketID)
	if s == nil {
		return fmt.Errorf("invalid socket")
	}
	if !*s.isConnecting {
		return fmt.Errorf("closed socket")
	}
	s.sendMsgChan <- msg
	return nil
}

func (server WsServer) BatchSend(msg []byte, ids ...int) {
	smap := server.sockets.Interface().(socketsMap)
	for _, id := range ids {
		go func(id int) {
			s := smap[id]
			if s != nil && *s.isConnecting {
				s.sendMsgChan <- msg
			}
		}(id)
	}
}

func (server WsServer) SendToAll(msg []byte) {
	smap := server.sockets.Interface().(socketsMap)
	for _, v := range smap {
		go func(v *socket) {
			if v != nil && *v.isConnecting {
				v.sendMsgChan <- msg
			}
		}(v)
	}
}

func (server WsServer) SendRoomMsg(msg []byte, roomID int) {
	server.roomMsg <- map[int][]byte{roomID: msg}
}
