package wsclient

import (
	"sync"

	"golang.org/x/net/websocket"
)

type WsClient struct {
	server            string
	origin            string
	conn              *websocket.Conn
	shutdownSignal    chan struct{}
	startSignal       chan struct{}
	stopSignal        *chan struct{}
	stoppedSignal     chan struct{}
	stopGroup         sync.WaitGroup
	isRunning         bool
	readChan          chan []byte
	writeChan         chan []byte
	mustReconnect     bool
	readErrorHandler  func(error)
	writeErrorHandler func(error)
	DebugOn           bool
}
