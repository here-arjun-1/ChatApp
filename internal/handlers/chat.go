package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chatapp/internal/store"
)

type ChatHandler struct {
	store *store.Store
}

func NewChatHandler(s *store.Store) *ChatHandler {
	return &ChatHandler{store: s}
}

type sendMessageRequest struct {
	From    string `json:"from" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := h.store.Add(req.From, req.Content)
	c.JSON(http.StatusCreated, msg)
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.History())
}
