package wsserver

import (
	"fmt"
	"net"
	"net/url"
	"testing"

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
	server.SetConnectionRouter(&addr, "/ws/test/")
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

	for i := 0; i < 100; i++ {
		client := wsclient.NewWsClient(s, origin, false, printError, printError)
		client.DebugOn = true
		client.Start()
		client.Send([]byte(fmt.Sprintf("ping %d", i)))
	}

	server.SendToAll([]byte("pong"))
	<-make(chan struct{})
}
