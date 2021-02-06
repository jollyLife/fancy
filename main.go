package main

import (
	"fancy/websocket"
	"fmt"
	"net/http"
)

func main() {
	s := websocket.New()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.Upgrade(w, r)
		fmt.Println(s.Get())
	})
	http.ListenAndServe(":8080", nil)
}
