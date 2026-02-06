package routes

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/xstar97/amialive/internal/api"
	"github.com/xstar97/amialive/internal/config"
	"net/http/httptest"
)

// Register sets up HTTP routes
func Register(cfg *config.Config) {
	// --- /healthz route (fun, random) ---
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		rand.Seed(time.Now().UnixNano())

		if rand.Intn(100) < cfg.JokeChance {
			fmt.Fprintln(w, api.GetJoke(cfg))
		} else {
			fmt.Fprintln(w, "pong!")
		}
	})

	// --- /healthcheck route (deterministic, machine-facing) ---
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		// Call /healthz handler directly without HTTP
		req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		rr := httptest.NewRecorder()

		// Use the same handler as /healthz
		http.DefaultServeMux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			http.Error(w, "healthz failed", http.StatusServiceUnavailable)
			return
		}

		// Everything good
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
}
