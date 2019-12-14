package main

import (
	"goutil/tcp"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := tcp.New(":9999")
	server.OnNewClient(func(c *tcp.Client) {
		println(`New client enter`)
	})
	server.OnNewMessage(func(c *tcp.Client, message string) {
		println(message)
	})
	server.OnClientConnectionClosed(func(c *tcp.Client, err error) {
		println(`Client closed`)
	})
	go server.Listen()
	time.Sleep(10 * time.Second)

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		println("Failed to connect to test server")
	}
	conn.Write([]byte("Test message\n"))
	conn.Close()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
