package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"chatapp/internal/store"
	"chatapp/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	store *store.Store
	hub   *ws.Hub
}

func NewWSHandler(s *store.Store, h *ws.Hub) *WSHandler {
	return &WSHandler{store: s, hub: h}
}

// incomingMessage covers every WS event a client can send:
//   - {"type":"message","content":"hi"}       -> new chat message
//   - {"type":"typing"} / {"type":"stop_typing"} -> typing indicator
//   - {"type":"read","ids":[1,2,3]}           -> mark messages seen
//
// Type is optional and defaults to "message" so older clients keep working.
type incomingMessage struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	IDs     []int  `json:"ids,omitempty"`
}

func (h *WSHandler) HandleConnection(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user query param required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	h.hub.Register(user, conn)
	defer h.hub.Unregister(user)

	for {
		var in incomingMessage
		if err := conn.ReadJSON(&in); err != nil {
			log.Println("read error:", err)
			break
		}

		switch in.Type {
		case "typing", "stop_typing":
			h.hub.Broadcast(user, gin.H{"type": in.Type, "from": user})

		case "read":
			if len(in.IDs) == 0 {
				continue
			}
			updated := h.store.MarkRead(in.IDs)
			if len(updated) == 0 {
				continue
			}
			h.hub.Broadcast(user, gin.H{"type": "read", "ids": updated, "by": user})

		default: // "message" or empty (older clients)
			if in.Content == "" {
				continue
			}
			msg := h.store.Add(user, in.Content)
			envelope := gin.H{"type": "message", "data": msg}
			h.hub.Broadcast(user, envelope)
			h.hub.SendTo(user, envelope)
		}
	}
}
