package main

import "fancy/server"

func main() {
	s := server.New()
	s.Listen()
}
