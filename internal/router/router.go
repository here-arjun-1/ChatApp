package router

import (
	"github.com/gin-gonic/gin"

	"chatapp/internal/handlers"
	"chatapp/internal/store"
)

func New(s *store.Store) *gin.Engine {
	r := gin.Default()

	chatHandler := handlers.NewChatHandler(s)

	r.GET("/health", handlers.Health)
	r.POST("/messages", chatHandler.SendMessage)
	r.GET("/messages", chatHandler.GetHistory)

	return r
}
