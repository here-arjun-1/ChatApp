package main

import (
	"chatapp/internal/router"
	"chatapp/internal/store"
)

func main() {
	s := store.New()
	r := router.New(s)
	r.Run(":9090")
}
