package wsclient

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var exitChan chan struct{}

func printError(err error) {
	// fmt.Printf("error: %v\n", err)
}

func TestClient(t *testing.T) {
	flag.Parse()
	server := "ws://121.40.165.18:8800"
	origin := "http://121.40.165.18:8800"
	var n int
	client := NewWsClient(server,
		WithOrigin(origin),
		WithReconnect(true),
		WithRecvErrorHandler(printError),
		WithSendErrorHandler(printError),
		WithReconnectingFunc(func() {
			n++
			fmt.Printf("n: %d\n", n)
		}),
		WithDebug(true),
	)

	client.Start()
	go func() {
		for {
			client.Receive()
		}
	}()
	time.Sleep(time.Second * 5)
	client.Close()
	time.Sleep(time.Second * 5)
}
