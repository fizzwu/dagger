package main

import (
	"log"
	"net"

	"time"

	"fmt"

	"github.com/fizzwu/dagger"
	"github.com/fizzwu/dagger/example/telnet"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":5555")
	if err != nil {
		log.Fatal("resolve tcp addr error:", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("listen err:", err)
	}

	server := dagger.NewServer(&telnet.TelnetCallback{}, &telnet.TelnetProtocol{}, 10, 10)

	fmt.Println("listening:", listener.Addr())
	server.Serve(listener, time.Second)

}
