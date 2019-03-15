package wsclient

import (
	"github.com/achillesss/go-utils/functionCaller"
	"github.com/achillesss/go-utils/log"
	"golang.org/x/net/websocket"
)

func (client *WsClient) debugLog(format string, args ...interface{}) {
	if client.DebugOn {
		log.Infofln(format, args...)
	}
}
func (client *WsClient) renewStopSignal() {
	sc := make(chan struct{})
	if client.stopSignal == nil {
		client.stopSignal = &sc
	} else {
		*client.stopSignal = sc
	}
}

// Connect establish a new websocket connection
// panic if any error occurs
func (client *WsClient) connect() {
	client.debugLog("start to connect")
	config, err := websocket.NewConfig(client.server, client.origin)
	if err != nil {
		panic(err)
	}
	conn, err := websocket.DialConfig(config)
	if err != nil {
		panic(err)
	}

	if client.conn != nil {
		client.conn.Close()
	}

	client.conn = conn
	client.isRunning = true
	client.debugLog("connected")
}

// Close close the connection in WsClient
// return false only when the closing method returns error
func (client *WsClient) Close() bool {
	if client.conn != nil {
		return client.conn.Close() == nil
	}
	return true
}

func (client *WsClient) stop() {
	client.debugLog("stop")
	close(*client.stopSignal)
	client.stopGroup.Wait()
	client.isRunning = false
	client.debugLog("stopped")
	client.renewStopSignal()
	if client.mustReconnect {
		client.stoppedSignal <- struct{}{}
	}
}

// 以下写法避免因读取消息时卡住而无法关闭连接
func (client *WsClient) read() {
	client.debugLog("start read")
	client.stopGroup.Add(1)
	defer client.stopGroup.Done()
	msgChan := make(chan string)
	go func() {
		for {
			var message string
			err := websocket.Message.Receive(client.conn, &message)
			if err != nil {
				if client.readErrorHandler != nil {
					client.readErrorHandler(err)
				}
				client.stop()
				return
			}
			msgChan <- message
		}
	}()

	for {
		select {
		case <-*client.stopSignal:
			client.debugLog("read routine stop")
			return
		case msg := <-msgChan:
			client.debugLog("reveice msg: %s", msg)
			client.readChan <- []byte(msg)
		}
	}
}

func (client *WsClient) write() {
	client.debugLog("start write")
	client.stopGroup.Add(1)
	defer client.stopGroup.Done()
	for {
		select {
		case <-*client.stopSignal:
			client.debugLog("write routine stop")
			return

		case msg := <-client.writeChan:
			client.debugLog("write msg: %s", msg)
			err := websocket.Message.Send(client.conn, string(msg))
			if err != nil {
				if client.writeErrorHandler != nil {
					client.writeErrorHandler(err)
				}
				client.stop()
			}
		}
	}
}

func (client *WsClient) start() {
	client.startSignal <- struct{}{}
}

func (client *WsClient) runMonitor() {
	client.debugLog("run monitor")
	go func() {
		for {
			select {
			case <-client.stoppedSignal:
				client.debugLog("receive stopped signal")
				go client.start()

			case <-client.shutdownSignal:
				client.debugLog("receive shutdown signal")
				client.stop()
				client.Close()
				return

				// start read and send
			case <-client.startSignal:
				client.debugLog("receive start signal")
				if client.isRunning {
					continue
				}
				client.connect()
				if client.connectingCount > 0 {
					client.connectingCompletedSignal <- struct{}{}
				}
				go client.read()
				go client.write()
				client.connectingCount++
			}
		}
	}()
}

var onConnectingCaller *funcaller.FunctionCaller

func SetOnReconnectingFunction(function interface{}, params ...interface{}) {
	onConnectingCaller = funcaller.NewCaller(function, params...)
}

func (client *WsClient) onReconnecting() {
	for range client.connectingCompletedSignal {
		onConnectingCaller.Call(false)
	}
}

func NewWsClient(server, origin string, mustReconnect bool, recvMsgErrHandler, sendMsgErrHandler func(error)) *WsClient {
	var client WsClient
	client.startSignal = make(chan struct{})
	client.shutdownSignal = make(chan struct{})
	client.stoppedSignal = make(chan struct{})
	client.readChan = make(chan []byte)
	client.writeChan = make(chan []byte)
	client.connectingCompletedSignal = make(chan struct{}, 1)
	client.server = server
	client.origin = origin
	client.readErrorHandler = recvMsgErrHandler
	client.writeErrorHandler = sendMsgErrHandler
	client.mustReconnect = mustReconnect
	client.renewStopSignal()

	return &client
}

func (client *WsClient) Start() {
	go client.onReconnecting()
	client.runMonitor()
	client.start()
}

func (client *WsClient) Send(msg []byte) {
	client.writeChan <- msg
}

func (client *WsClient) Receive() []byte {
	return <-client.readChan
}

func (client *WsClient) Reconnect() {
	client.stop()
	client.Close()
	client.start()
}
