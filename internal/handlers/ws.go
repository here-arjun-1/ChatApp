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

type incomingMessage struct {
	Content string `json:"content"`
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

		msg := h.store.Add(user, in.Content)
		h.hub.Broadcast(user, msg)
	}
}
