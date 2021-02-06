package main

import (
	"fmt"
	"net/url"
)

func main() {

	a := `tcp@:3109:calls=0.00&connections=0.00&%7B%22server_id%22%3A%225c5448d3-9cb2-4ce4-866e-cabe011e4307%22%2C%22ip%22%3A%22%3A3109%22%7D=]`

	d, err := url.PathUnescape(a)
	if err != nil {
		return
	}

	fmt.Println(d)

	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//
	//	fmt.Println(r.Header)
	//
	//	server.Websocket()
	//})
	//
	//http.ListenAndServe(":8080", nil)
}
