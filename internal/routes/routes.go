package routes

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/xstar97/amialive/internal/api"
	"github.com/xstar97/amialive/internal/config"
)

// Register sets up HTTP routes
func Register(cfg *config.Config) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		rand.Seed(time.Now().UnixNano())

		if rand.Intn(100) < cfg.JokeChance {
			fmt.Fprintln(w, api.GetJoke())
		} else {
			fmt.Fprintln(w, "pong!")
		}
	})
}
