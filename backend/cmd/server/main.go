package main

import (
	"os"

	"chatapp/internal/router"
	"chatapp/internal/store"
	"chatapp/internal/ws"
)

func main() {
	s := store.New()
	hub := ws.NewHub()
	r := router.New(s, hub)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	r.Run(":" + port)
}
