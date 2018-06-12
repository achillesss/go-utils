package wsserver

import (
	"github.com/achillesss/go-utils/go-map"
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

type WsServer struct {
	handlers            []*wsHandler
	sockets             *gomap.GoMap
	lastConnID          *int
	receiveMsgHandler   func([]byte)
	receiveErrorHandler func(error)
	sendErrorHandler    func(error)
	ShouldAnswerPing    bool // hot switch
}
