package wsclient

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var exitChan chan struct{}

func printError(err error) {
	fmt.Printf("error: %v\n", err)
}

func TestClient(t *testing.T) {
	flag.Parse()
	server := "ws://118.25.40.163:8088"
	origin := "http://118.25.40.163"
	client := NewWsClient(server, origin, true, printError, printError)
	client.DebugOn = true
	client.Start()
	go func() {
		for {
			msg := client.Receive()
			fmt.Printf("receive %s", msg)
		}
	}()
	time.Sleep(time.Second * 5)
	client.stop()
	time.Sleep(time.Second * 5)
}
