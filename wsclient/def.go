package wsclient

import (
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type WsClient struct {
	server                    string
	origin                    string
	header                    http.Header
	conn                      *websocket.Conn
	shutdownSignal            chan struct{}
	startSignal               chan struct{}
	connectingCompletedSignal chan struct{}
	stopSignal                *chan struct{}
	stoppedSignal             chan struct{}
	stopGroup                 sync.WaitGroup
	isRunning                 bool
	readChan                  chan []byte
	writeChan                 chan []byte
	mustReconnect             bool
	readErrorHandler          func(error)
	writeErrorHandler         func(error)
	DebugOn                   bool
	connectingCount           int
	onReconnectingFunc        func()
}
