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
var client = NewWsClient(
    wsAddr, // ws server address
    WithOrigin(origin), // if connect with origin
    WithHeader(header), // if connect with header
    WithReconnect(reconnectIfDisconnected), // if want to reconnect when socket disconnected
    WithRecvErrorHandler(handleRecvMsgErrFunc), // to handle error on receiving msg
    WithSendErrorHandler(handleSendMsgErrFunc), // to handle error on sending msg
    WithDebug(true), // turn on debug logs
    WithReconnectingFunc(func(){ print("reconnected!") }), // run function on reconnecting
)

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

