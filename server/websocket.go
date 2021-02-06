package server

import (
	"fmt"
	"net"
)

func Websocket() {
	fmt.Println("ws")

	l, err := net.Listen("tcp", ":3102")
	if err != nil {
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			return
		}

		c.LocalAddr()
	}

}
