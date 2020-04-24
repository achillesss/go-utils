package wsserver

import (
	"net/http"
	"sync"
)

type MsgHandler func([]byte)
type ErrHandler func(error)

type WsServer struct {
	lock    sync.Mutex
	routers map[string]map[string]http.Handler

	// rooms               *gomap.GoMap
	lastConnID int64
	// lastRoomID          int64
	receiveMsgHandler   MsgHandler
	receiveErrorHandler ErrHandler
	sendMsgHandler      MsgHandler
	sendErrorHandler    ErrHandler
	ShouldAnswerPing    bool // hot switch
	// roomMsg             chan map[int64][]byte

	socketMsgChanLock sync.RWMutex
	socketMsgChans    map[int64]chan []byte

	socketStopChanLock sync.RWMutex
	socketStopChans    map[int64]chan struct{}
}

type ServerOption interface {
	updateOption(*WsServer)
}

type optionHolder struct {
	updateFunc func(*WsServer)
}
