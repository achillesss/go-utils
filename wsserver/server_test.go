package wsserver

import (
	"fmt"
	"net"
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
	var addr net.TCPAddr
	addr.Zone = ""
	addr.Port = 9999
	addr.IP = net.IPv4(127, 0, 0, 1)

	config := NewEmptyConfig()
	config.SetReceiveErrorHandler(printError)
	config.SetSendErrorHandler(printError)
	config.SetReceiveMsgHandler(printMsg)
	server := NewWsServerFromConfig(config)
	server.SetConnectionRouters(&addr, "/ws/test/")
	server.Serve()
	server.ShouldAnswerPing = true

	var u url.URL
	u.Scheme = "ws"
	u.Host = addr.String()
	u.Path = "/ws/test/"
	s := u.String()
	u.Scheme = "http"
	u.Path = ""
	origin := u.String()

	for i := 0; i < 3; i++ {
		client := wsclient.NewWsClient(s, origin, false, printError, printError)
		client.DebugOn = true
		client.Start()
		client.Send([]byte(fmt.Sprintf("ping %d", i)))
		go func() {
			for {
				client.Receive()
			}
		}()
	}

	server.SendToAll([]byte("pong"))

	time.Sleep(time.Second)

	server.NewRoom()
	server.NewRoom()
	server.NewRoom()
	server.NewRoom()
	server.NewRoom()

	roomIDs := server.ListRooms()
	fmt.Printf("roomIDs: %d\n", roomIDs)

	socketIDs := server.ListSockets()
	fmt.Printf("sockets: %d\n", socketIDs)

	r := server.queryRoom(roomIDs[0])
	server.AddRoomate(roomIDs[0], socketIDs[0])
	roomates := r.ListRoomates()
	fmt.Printf("roomates: %d\n", roomates)
	server.AddRoomate(roomIDs[0], socketIDs[1])
	roomates = r.ListRoomates()
	fmt.Printf("roomates: %d\n", roomates)
	server.AddRoomate(roomIDs[0], socketIDs[2])
	roomates = r.ListRoomates()
	fmt.Printf("roomates: %d\n", roomates)

	server.SendRoomMsg([]byte(fmt.Sprintf("welcome to room %d", roomIDs[0])), roomIDs[0])
	<-make(chan struct{})
}
