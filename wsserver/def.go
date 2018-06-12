package wsserver

import (
	"github.com/achillesss/go-utils/go-map"
)

type WsServer struct {
	handlers            []*wsHandler
	sockets             *gomap.GoMap
	rooms               *gomap.GoMap
	lastConnID          *int
	lastRoomID          *int
	receiveMsgHandler   func([]byte)
	receiveErrorHandler func(error)
	sendErrorHandler    func(error)
	ShouldAnswerPing    bool // hot switch
	roomMsg             chan map[int][]byte
}
