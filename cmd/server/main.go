package main

import (
	"chatapp/internal/router"
	"chatapp/internal/store"
	"chatapp/internal/ws"
)

func main() {
	s := store.New()
	hub := ws.NewHub()
	r := router.New(s, hub)
	r.Run(":9090")
}
