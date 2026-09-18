package router
<<<<<<< HEAD
=======

import (
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.Health)

	return r
}
>>>>>>> f0c79c5 (feat: add health router)
