package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"time"
)

// Serve 是一个对外的接口，需要实现此接口进行消息处理
type Serve interface {
	Conn(net.Conn)
}

// Server 结构体，主要包含了tcp服务器的一些属性
type Server struct{}

// New 创建一个tcp server 持续监听连接
func New() *Server {
	s := &Server{}
	return s
}

// Listen 监听客户端发来的连接请求
func (s *Server) Listen(addr ...string) {
	address := resolveAddress(addr)
	s.listen(address)
}

// 监听
func (s *Server) listen(addr string) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("tcp服务器已启动...正在监听端口:", addr)
	s.accept(listen)
}

func (s *Server) accept(listen net.Listener) {
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		go s.Conn(conn)
	}
}

func (s *Server) Conn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到client发来的数据：", msg)

		ch := make(chan struct{}, 1)
		ch <- struct{}{}
		heartbeat(conn, ch)

		send, _ := Encode(fmt.Sprintf("goroutine, nums:%d, client: %s", runtime.NumGoroutine(), strings.Split(msg, ":")[1]))
		conn.Write(send)
	}
}

func heartbeat(conn net.Conn, ping chan struct{}) {
	select {
	case <-ping:
		fmt.Println("keep breathing")
		conn.SetDeadline(time.Now().Add(time.Second * 5))
	case <-time.After(time.Second * 5):
		fmt.Println("--------------")
		conn.Close()
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		return ":30000"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
