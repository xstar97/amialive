package main

import (
	"fmt"
	"net/http"

	"github.com/xstar97/amialive/internal/config"
	"github.com/xstar97/amialive/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	routes.Register(cfg)

	fmt.Printf("Server running on port %s (JOKE_CHANCE=%d%%)\n", cfg.Port, cfg.JokeChance)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		panic(err)
	}
}
