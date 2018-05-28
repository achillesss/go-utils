package wsclient

import (
	"time"

	"bitbucket.org/magmeng/go-utils/log"
	"golang.org/x/net/websocket"
)

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
	log.Infofln("start to connect")
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
	log.Infofln("connected")
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
	log.Infofln("stop")
	close(*client.stopSignal)
	client.stopGroup.Wait()
	client.isRunning = false
	time.Sleep(time.Second)
	log.Infofln("stopped")
	client.renewStopSignal()
	if client.mustReconnect {
		client.stoppedSignal <- struct{}{}
	}
}

// 以下写法避免因读取消息时卡住而无法关闭连接
func (client *WsClient) read() {
	log.Infofln("start read")
	client.stopGroup.Add(1)
	defer client.stopGroup.Done()
	msgChan := make(chan string)

	go func() {
		for {
			var message string
			err := websocket.Message.Receive(client.conn, &message)
			if err != nil && client.readErrorHandler != nil {
				client.readErrorHandler(err)
				client.stop()
				return
			}
			msgChan <- message
		}
	}()

	for {
		select {
		case <-*client.stopSignal:
			log.Infofln("read routine stop")
			return
		case msg := <-msgChan:
			log.Infofln("reveice msg: %s", msg)
			client.readChan <- []byte(msg)
		}
	}
}

func (client *WsClient) write() {
	log.Infofln("start write")
	client.stopGroup.Add(1)
	defer client.stopGroup.Done()
	for {
		select {
		case <-*client.stopSignal:
			log.Infofln("write routine stop")
			return
		case msg := <-client.writeChan:
			log.Infofln("write msg: %s", msg)
			err := websocket.Message.Send(client.conn, string(msg))
			if err != nil && client.writeErrorHandler != nil {
				client.writeErrorHandler(err)
				client.stop()
			}
		}
	}
}

func (client *WsClient) start() {
	client.startSignal <- struct{}{}
}

func (client *WsClient) reconnect() {
	client.stop()
	client.Close()
	client.start()
}

func (client *WsClient) runMonitor() {
	log.Infofln("run monitor")

	go func() {
		for {
			select {
			case <-client.stoppedSignal:
				log.Infofln("receive stopped signal")
				go client.start()
			case <-client.shutdownSignal:
				log.Infofln("receive shutdown signal")
				client.stop()
				client.Close()
				return

			// start read and send
			case <-client.startSignal:
				log.Infofln("receive start signal")
				if client.isRunning {
					continue
				}
				client.connect()
				go client.read()
				go client.write()
			}
		}
	}()

}

func NewWsClient(server, origin string, mustReconnect bool, recvMsgErrHandler, sendMsgErrHandler func(error)) *WsClient {
	log.Infofln("init new client")
	var client WsClient
	client.startSignal = make(chan struct{})
	client.shutdownSignal = make(chan struct{})
	client.stoppedSignal = make(chan struct{})
	client.readChan = make(chan []byte)
	client.writeChan = make(chan []byte)
	client.server = server
	client.origin = origin
	client.readErrorHandler = recvMsgErrHandler
	client.writeErrorHandler = sendMsgErrHandler
	client.mustReconnect = mustReconnect
	client.renewStopSignal()

	client.runMonitor()
	client.start()
	return &client

}

func (client *WsClient) Send(msg []byte) {
	client.writeChan <- msg
}

func (client *WsClient) Receive() []byte {
	return <-client.readChan
}
