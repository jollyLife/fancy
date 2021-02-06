package websocket

import (
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	ErrClientNotFound = errors.New("client not found")
)

// Server
type Server struct {
	sync.RWMutex
	size *int64
	// 存放每个客户端的信息
	client map[int64]*Client
}

// New
func New() *Server {
	s := &Server{
		size:   new(int64),
		client: make(map[int64]*Client, 128),
	}
	return s
}

func (s *Server) register(client *Client) bool {
	s.Lock()
	defer s.Unlock()

	// 检查有没有注册
	_, ok := s.client[client.id]
	if ok {
		return false
	}

	s.client[client.id] = client

	return true
}

func unregister(clientId int64, server *Server) {
	server.Lock()
	defer server.Unlock()
	delete(server.client, clientId)
}

// GetClientById
func GetClientById(clientId int64, server *Server) (*Client, error) {
	server.RLock()
	defer server.RUnlock()

	client, ok := server.client[clientId]
	if !ok {
		return nil, ErrClientNotFound
	}

	return client, nil
}

// Upgrade
func (s *Server) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientId := atomic.AddInt64(s.size, 1)
	client := &Client{
		id:   clientId,
		conn: conn,
		send: make(chan []byte, 256),
	}

	s.register(client)

	go client.read(s)
	go client.write()
}

func (s *Server) Get() map[int64]*Client {
	return s.client
}
