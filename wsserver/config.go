package wsserver

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func ConvertHttpHandlerFunctionToHttpHandler(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.Handler(http.HandlerFunc(f))
}

func ConvertWsHandlerFunctionToHttpHandler(f func(c *websocket.Conn)) http.Handler {
	return websocket.Handler(f)
}

func (o *optionHolder) updateOption(ws *WsServer) {
	o.updateFunc(ws)
}

func newOptionHolder(f func(ws *WsServer)) *optionHolder {
	return &optionHolder{
		updateFunc: f,
	}
}
