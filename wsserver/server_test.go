package wsserver

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/achillesss/go-utils/wsclient"
)

func printError(err error) {
	println(fmt.Sprintf("error: %v\n\n", err))
}
func printMsg(msg []byte) {
	println(fmt.Sprintf("msg: %s\n\n", msg))

}

func TestServer(t *testing.T) {
	server := NewWsServer(
		"127.0.0.1:9999", "/ws/test/",
		WithRecvErrorHandler(printError),
		WithSendErrorHandler(printError),
		WithRecvMsgHandler(func(m []byte) { printMsg(m) }),
		WithAnswerPing(),
	)
	server.Serve()

	var u url.URL
	u.Scheme = "ws"
	u.Host = "127.0.0.1:9999"
	u.Path = "/ws/test/"
	s := u.String()
	u.Scheme = "http"
	u.Path = ""
	origin := u.String()

	for i := 0; i < 3; i++ {
		client := wsclient.NewWsClient(s, origin, nil, false, printError, printError)
		client.DebugOn = true
		client.Start()
		client.Send([]byte(fmt.Sprintf("ping %d", i)))
		go func() {
			for {
				msg := client.Receive()
				fmt.Printf("client receive: %s\n", msg)
			}
		}()
	}

	server.Broadcast([]byte("pong"))
	server.SendMsgTo([]byte("ppp"), 1, 2)

	time.Sleep(time.Second)

	//	server.NewRoom()
	//	server.NewRoom()
	//	server.NewRoom()
	//	server.NewRoom()
	//	server.NewRoom()
	//
	//	roomIDs := server.ListRooms()
	//	fmt.Printf("roomIDs: %d\n", roomIDs)
	//
	socketIDs := server.ListSocketIDs()
	fmt.Printf("sockets: %d\n", socketIDs)
	//
	//	r := server.queryRoom(roomIDs[0])
	//	server.AddRoomates(roomIDs[0], socketIDs[0])
	//	roomates := r.ListRoomates()
	//	fmt.Printf("roomates: %d\n", roomates)
	//	server.AddRoomates(roomIDs[0], socketIDs[1])
	//	roomates = r.ListRoomates()
	//	fmt.Printf("roomates: %d\n", roomates)
	//	server.AddRoomates(roomIDs[0], socketIDs[2])
	//	roomates = r.ListRoomates()
	//	fmt.Printf("roomates: %d\n", roomates)
	//
	//	server.SendRoomMsg([]byte(fmt.Sprintf("welcome to room %d", roomIDs[0])), roomIDs[0])
	<-make(chan struct{})
}
