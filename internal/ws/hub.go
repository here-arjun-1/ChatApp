package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu    sync.RWMutex
	conns map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		conns: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) Register(user string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[user] = conn
}

func (h *Hub) Unregister(user string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns, user)
}

func (h *Hub) Broadcast(sender string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for user, conn := range h.conns {
		if user == sender {
			continue
		}
		conn.WriteJSON(data)
	}
}
