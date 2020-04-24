package wsclient

import "net/http"

type ClientOption interface {
	updateOption(*WsClient)
}

type optionHolder struct {
	updateFunc func(*WsClient)
}

func (h *optionHolder) updateOption(c *WsClient) {
	h.updateFunc(c)
}

func newOptionHolder(f func(*WsClient)) ClientOption {
	return &optionHolder{
		updateFunc: f,
	}
}

func WithOrigin(origin string) ClientOption {
	return newOptionHolder(func(c *WsClient) {
		c.origin = origin
	})
}

func WithHeader(header http.Header) ClientOption {
	return newOptionHolder(func(c *WsClient) {
		c.header = header
	})
}

func WithReconnect(mustReconnect bool) ClientOption {
	return newOptionHolder(func(c *WsClient) {
		c.mustReconnect = mustReconnect
	})
}

func WithRecvErrorHandler(f func(error)) ClientOption {
	return newOptionHolder(func(c *WsClient) {
		c.readErrorHandler = f
	})
}

func WithSendErrorHandler(f func(error)) ClientOption {
	return newOptionHolder(func(c *WsClient) {
		c.writeErrorHandler = f
	})
}

func WithReconnectingFunc(f func()) ClientOption {
	return newOptionHolder(func(c *WsClient) {
		c.onReconnectingFunc = f
	})
}

func WithDebug(on bool) ClientOption {
	return newOptionHolder(func(c *WsClient) {
		c.DebugOn = on
	})
}
