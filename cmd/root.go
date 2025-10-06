package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xstar97/amialive/internal/config"
	"github.com/xstar97/amialive/internal/routes"
)

func main() {
	cfg := config.Load()

	routes.Register(cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
