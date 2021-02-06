package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
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
	}

	unregister(c.id, server)
}

func (c *Client) hub(message []byte) {
	GetClientById()
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
