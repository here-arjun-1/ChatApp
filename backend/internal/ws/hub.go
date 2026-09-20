package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// client wraps a connection with its own write lock. gorilla/websocket only
// allows one concurrent writer per connection, and once we broadcast to a
// user (typing, read receipts, messages) while also echoing directly back to
// that same user, writes can race without this.
type client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *client) writeJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(v)
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*client),
	}
}

func (h *Hub) Register(user string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[user] = &client{conn: conn}
}

func (h *Hub) Unregister(user string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, user)
}

// Broadcast sends data to every connected user except sender.
func (h *Hub) Broadcast(sender string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for user, c := range h.clients {
		if user == sender {
			continue
		}
		c.writeJSON(data)
	}
}

// SendTo sends data directly to a single user, if connected.
func (h *Hub) SendTo(user string, data interface{}) {
	h.mu.RLock()
	c, ok := h.clients[user]
	h.mu.RUnlock()
	if !ok {
		return
	}
	c.writeJSON(data)
}
