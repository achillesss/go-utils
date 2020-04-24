wsclient

websocket 客户端工具
使用方法:
```golang
var wsAddr = "ws://127.0.0.1:1234"
var origin = "http://127.0.0.1:1234"
var header = make(map[string]string)
var reconnectIfDisconnected = true
var handleRecvMsgErrFunc = func(err error){ print(err) }
var handleSendMsgErrFunc = func(err error){ print(err) }
var client = NewWsClient(wsAddr, origin, header, reconnectIfDisconnected, handleRecvMsgErrFunc, handleSendMsgErrFunc)

// if you want to do something on reconnecting:
client.SetOnReconnectingFunction(f func(params ...), params ...)
// if you want to debug
// client.DebugOn = true

// start client
client.Start()

// Receive msg
for {
    var msg = client.Receive()
    // do something
}

// Send msg
client.Send([]byte("hello, websocket!"))


// stop client
client.Close()
```        

