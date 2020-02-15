package main

import (
	"github.com/jzy/go-mmo-server/src/message"
	"log"
	"net"
	"net/rpc"
)

type HelloService struct {
}

func (p *HelloService) Hello(request message.String, response *message.String) error {
	response.Value = "hello:" + request.GetValue()
	return nil
}

func main() {
	rpc.Register(new(HelloService))
	listener, error := net.Listen("tcp", ":1234")
	if error != nil {
		log.Fatal("ListenTcp error:", error)
	}
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go rpc.ServeConn(con)
	}
}
