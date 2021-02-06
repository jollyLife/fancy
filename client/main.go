package main

import (
	"fancy/server"
	"fmt"
	"net"
	"time"
)

func dial(num int) {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	msg := fmt.Sprintf(`$msg $1 $message num:%d`, num)
	data, err := server.Encode(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	//reader := bufio.NewReader(conn)

	for {
		_, err := conn.Write(data)
		fmt.Println(err)
		time.Sleep(time.Second * 5)
		fmt.Println("goroutine:")
		//s, err := server.Decode(reader)
		//if err != nil {
		//	fmt.Println(err)
		//	continue
		//}
		//fmt.Println("goroutine:", s)
	}

}

func main() {
	for i := 0; i < 1; i++ {
		go dial(i)
	}

	select {}
}
