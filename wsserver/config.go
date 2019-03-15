package wsserver

import (
	"net"
	"net/http"

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
	handlers []*wsHandler
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

func NewEmptyConfig() *wsServerConfig {
	var c wsServerConfig
	return &c
}
