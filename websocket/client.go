package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client
type Client struct {
	id   int64
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read(server *Server) {
	defer func() {
		c.conn.Close()
	}()
	c.send <- []byte(fmt.Sprintf("$1$%d", c.id))
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(message)
		hub(server, message)
	}
	unregister(c.id, server)
}

func hub(server *Server, message []byte) {
	msg := string(message)
	id := strings.Split(msg, "$")[1]
	otherId, _ := strconv.ParseInt(id, 10, 64)
	other, err := GetClientById(otherId, server)
	if err != nil {
		fmt.Println(err)
		return
	}

	other.send <- message
}

func (c *Client) write() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)

			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
