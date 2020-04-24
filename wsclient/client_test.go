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
	client := NewWsClient(server, origin, nil, true, printError, printError)
	var n int
	SetOnReconnectingFunction(func(n *int) {
		*n++
		fmt.Printf("n: %d\n", *n)
	}, &n)

	client.DebugOn = true
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
