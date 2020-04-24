wsserver
websocket 服务端工具
使用方法:
```golang
var addr = "127.0.0.1:1234"
var pattern = "/ws/test/"
var anotherHandler http.Handler

// create new server
var server = NewWsServer(
    addr, pattern,
    WithSendErrorHandler(func(err error){ print(err) }), // if you want to handle error on sending msg
    WithSendMsgHandler(func(msg []byte){ print(msg) }), // if you want to handle msg on sending msg
    WithRecvErrorHandler(func(err error){ print(err) }), // if you want to handle error on receiving msg
    WithRecvMsgHandler(func(msg []byte){ print(msg) }), // if you want to handle msg on receiving msg
    WithRouter("127.0.0.1:1234", "/ws/test2/", anotherhandler), // if you want to add other handlers
    WithAnswerPing(), // if you want to reponse with 'pong' when receiving 'ping'
)

// serve
server.Serve()


// list known socket ids:
var socketIDs = server.ListSocketIDs()

// send msg to specific sockets
server.SendMsgTo([]byte("hello, socket!"), socketIDs...)

// broadcast msg
server.Broadcast([]byte("hello, sockets!"))

// disconnect to a socket:
server.Stop(socketIDs[0])
```
