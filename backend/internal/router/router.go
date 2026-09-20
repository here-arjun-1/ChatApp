package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"chatapp/internal/handlers"
	"chatapp/internal/store"
	"chatapp/internal/ws"
)

func corsMiddleware() gin.HandlerFunc {
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func New(s *store.Store, hub *ws.Hub) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	chatHandler := handlers.NewChatHandler(s)
	wsHandler := handlers.NewWSHandler(s, hub)

	r.GET("/health", handlers.Health)
	r.POST("/messages", chatHandler.SendMessage)
	r.GET("/messages", chatHandler.GetHistory)
	r.GET("/ws", wsHandler.HandleConnection)

	return r
}
