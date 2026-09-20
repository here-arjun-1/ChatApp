package router

import (
	"github.com/gin-gonic/gin"

	"chatapp/internal/handlers"
	"chatapp/internal/store"
	"chatapp/internal/ws"
)

func New(s *store.Store, hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	chatHandler := handlers.NewChatHandler(s)
	wsHandler := handlers.NewWSHandler(s, hub)

	r.GET("/health", handlers.Health)
	r.POST("/messages", chatHandler.SendMessage)
	r.GET("/messages", chatHandler.GetHistory)
	r.GET("/ws", wsHandler.HandleConnection)

	return r
}
