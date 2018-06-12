package wsserver

import (
	"fmt"
	"net"
	"net/http"

	"github.com/achillesss/go-utils/go-map"
	"github.com/achillesss/go-utils/log"
	"golang.org/x/net/websocket"
)

type wsRouter struct {
	pattern string
	handler http.Handler
}

func ConvertHttpHandlerFunctionToHttpHandler(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.Handler(http.HandlerFunc(f))
}

func ConvertWsHandlerFunctionToHttpHandler(f func(c *websocket.Conn)) http.Handler {
	return websocket.Handler(f)
}

func NewWsRouter(pattern string, handler http.Handler) *wsRouter {
	var r wsRouter
	r.pattern = pattern
	r.handler = handler
	return &r
}

type wsRouters map[string]http.Handler

type wsHandler struct {
	address net.Addr
	routers wsRouters
}

func (h *wsHandler) registerRouter(pattern string, handler http.Handler) {
	h.routers[pattern] = handler
}

func (h wsHandler) serve() {
	l, err := net.Listen(h.address.Network(), h.address.String())
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	for k, v := range h.routers {
		log.Infofln("new router: %s", k)
		mux.Handle(k, v)
	}

	go http.Serve(l, mux)
	log.Infofln("serve above routers on %s", l.Addr().String())
}

type wsServerConfig struct {
	handlers         []*wsHandler
	recvErrorHandler func(error)
	recvMsgHandler   func([]byte)
	sendErrorHandler func(error)
}

func newWsHandler() *wsHandler {
	var h wsHandler
	h.routers = make(wsRouters)
	return &h
}

func (c *wsServerConfig) AddHandler(addr net.Addr, routers ...wsRouter) {
	h := newWsHandler()
	h.address = addr
	for _, r := range routers {
		h.registerRouter(r.pattern, r.handler)
	}
	c.handlers = append(c.handlers, h)
}

func (c *wsServerConfig) SetSendErrorHandler(f func(error)) {
	c.sendErrorHandler = f
}

func (c *wsServerConfig) SetReceiveErrorHandler(f func(error)) {
	c.recvErrorHandler = f
}

func (c *wsServerConfig) SetReceiveMsgHandler(f func([]byte)) {
	c.recvMsgHandler = f
}

func NewEmptyConfig() *wsServerConfig {
	var c wsServerConfig
	return &c
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

func (server *WsServer) initSocketsMap() {
	server.sockets = gomap.NewMap(make(socketsMap))
	go server.sockets.Handler()
}

func (server *WsServer) onConnection(c *websocket.Conn) {
	log.Infofln("new connection")
	s := server.newSocket(c)
	server.sockets.Add(s.id, s)
	go s.send(server.sendErrorHandler)
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

func NewWsServerFromConfig(config *wsServerConfig) *WsServer {
	var server WsServer
	server.lastConnID = new(int)
	server.handlers = config.handlers
	server.initSocketsMap()
	server.sendErrorHandler = config.sendErrorHandler
	server.receiveErrorHandler = config.recvErrorHandler

	server.receiveMsgHandler = config.recvMsgHandler
	return &server
}

func (server WsServer) SetConnectionRouter(addr net.Addr, pattern string) {
	h := newWsHandler()
	h.address = addr
	h.registerRouter(pattern, ConvertWsHandlerFunctionToHttpHandler(server.onConnection))
	h.serve()
}

func (server WsServer) Serve() {
	for _, h := range server.handlers {
		h.serve()
	}
}

func (server WsServer) SendTo(msg []byte, id int) error {
	var s *socket
	server.sockets.Query(id, &s)
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
